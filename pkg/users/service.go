package users

import (
	"context"
	"strconv"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/auth"
	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error)
	Login(ctx context.Context, input *model.UserInput, tmpUserID string) (*model.UserOutput, bool, error)
	RefreshToken(ctx context.Context, token string) (*model.UserOutput, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
}

type Storage interface {
	RegisterUser(ctx context.Context, user *model.User) (string, error)
	Login(ctx context.Context, email string, phone int64, tmpUserID string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
}

type BasicService struct {
	logger  log.Logger
	storage Storage
}

type ServiceConfig struct {
	Logger  log.Logger
	Storage Storage
}

func (s *BasicService) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {

	createdTime := time.Now()

	err := input.Validate()
	if err != nil {
		return nil, errBadRequest("validation error: %v", err)
	}

	newUser := input.Init(createdTime)
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 5)
	if err != nil {
		return nil, errInternal("%s", err)
	}

	newUser.Password = string(hash)
	newUser.ID = primitive.NewObjectID()
	userID, err := s.storage.RegisterUser(ctx, newUser)
	if err != nil {
		return nil, errUnauthorized("%s", err)
	}

	var claimSubject string
	if input.Email != "" {
		claimSubject = input.Email
	} else {
		claimSubject = strconv.FormatInt(input.Phone, 10)
	}

	//expirationTime := time.Now().Add(5 * time.Minute)
	claims := &auth.Claims{
		Subject:        claimSubject,
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(auth.JwtKey)
	if err != nil {
		return nil, errUnauthorized("%s", err)
	}

	userOutput := &model.UserOutput{}
	userOutput.ID = userID
	userOutput.Token = tokenString

	return userOutput, nil
}

func (s *BasicService) Login(ctx context.Context, input *model.UserInput, tmpUserID string) (*model.UserOutput, bool, error) {

	err := input.Validate()
	if err != nil {
		return nil, false, errBadRequest("validation error: %v", err)
	}

	result, err := s.storage.Login(ctx, input.Email, input.Phone, tmpUserID)
	if err != nil {
		return nil, false, errBadRequest("%s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(input.Password))
	if err != nil {
		return nil, false, errUnauthorized("%s", err)
	}

	var claimSubject string
	if input.Email != "" {
		claimSubject = input.Email
	} else {
		claimSubject = strconv.FormatInt(input.Phone, 10)
	}

	//expirationTime := time.Now().Add(5 * time.Minute)
	claims := &auth.Claims{
		Subject:        claimSubject,
		UserID:         result.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(auth.JwtKey)
	if err != nil {
		return nil, false, errInternal("%s", err)
	}

	userOutput := &model.UserOutput{}
	userOutput.ID = result.ID.Hex()
	userOutput.Token = tokenString

	var unsetTmpUserIDCookie bool
	if tmpUserID != "" {
		unsetTmpUserIDCookie = true
	}

	return userOutput, unsetTmpUserIDCookie, nil
}

func (s *BasicService) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
	claims := &auth.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})
	s.logger.Log("token", token)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errUnauthorized("%s", err)
		}
		return nil, errBadRequest("%s", err)
	}
	if !tkn.Valid {
		return nil, errUnauthorized("%s", err)
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return nil, errBadRequest("%s", err)
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(auth.JwtKey)
	if err != nil {
		return nil, errInternal("%s", err)
	}

	userOutput := &model.UserOutput{}
	userOutput.Token = tokenString
	userOutput.ExpirationTime = expirationTime

	return userOutput, nil
}

func (s *BasicService) GetUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		return nil, errInternal("%s", err)
	}
	return users, nil
}

func NewService(cfg ServiceConfig) (*BasicService, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = log.NewNopLogger()
	}

	if cfg.Storage == nil {
		return nil, errBadRequest("storage must be provided")
	}

	svc := &BasicService{
		logger:  logger,
		storage: cfg.Storage,
	}

	return svc, nil
}
