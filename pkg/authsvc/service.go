package erpsvc

import (
	"context"
	"github.com/anabiozz/core/lapkins/pkg/auth"
	"github.com/anabiozz/core/lapkins/pkg/erp"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type Storage interface {
	RegisterUser(ctx context.Context, user *erp.User) (string, error)
	Login(ctx context.Context, email string, phone int64, tmpUserID string) (*erp.User, error)
	GetUsers(ctx context.Context) ([]*erp.User, error)
}

type Service interface {
	Register(ctx context.Context, input *erp.UserInput) (*erp.UserOutput, error)
	Login(ctx context.Context, input *erp.UserInput, tmpUserID string) (*erp.UserOutput, bool, error)
	RefreshToken(ctx context.Context, token string) (*erp.UserOutput, error)
	GetUsers(ctx context.Context) ([]*erp.User, error)
}

type service struct {
	logger  log.Logger
	storage Storage
}

type ServiceConfig struct {
	Logger  log.Logger
	Storage Storage
}

func newService(cfg *ServiceConfig) (*service, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = log.NewNopLogger()
	}

	svc := &service{
		logger:  logger,
		storage: cfg.Storage,
	}

	return svc, nil
}

func (s *service) Register(ctx context.Context, input *erp.UserInput) (*erp.UserOutput, error) {
	createdTime := time.Now()

	err := input.Validate()
	if err != nil {
		return nil, erp.ErrBadRequest("validation error: %v", err)
	}

	newUser := input.Init(createdTime)
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 5)
	if err != nil {
		return nil, erp.ErrInternal("%s", err)
	}

	newUser.Password = string(hash)
	newUser.ID = primitive.NewObjectID()
	userID, err := s.storage.RegisterUser(ctx, newUser)
	if err != nil {
		return nil, erp.ErrUnauthorized("%s", err)
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
		return nil, erp.ErrUnauthorized("%s", err)
	}

	userOutput := &erp.UserOutput{}
	userOutput.ID = userID
	userOutput.Token = tokenString

	return userOutput, nil
}

func (s *service) Login(ctx context.Context, input *erp.UserInput, tmpUserID string) (*erp.UserOutput, bool, error) {
	err := input.Validate()
	if err != nil {
		return nil, false, erp.ErrBadRequest("validation error: %v", err)
	}

	result, err := s.storage.Login(ctx, input.Email, input.Phone, tmpUserID)
	if err != nil {
		return nil, false, erp.ErrBadRequest("%s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(input.Password))
	if err != nil {
		return nil, false, erp.ErrUnauthorized("%s", err)
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
		return nil, false, erp.ErrInternal("%s", err)
	}

	userOutput := &erp.UserOutput{}
	userOutput.ID = result.ID.Hex()
	userOutput.Token = tokenString

	var unsetTmpUserIDCookie bool
	if tmpUserID != "" {
		unsetTmpUserIDCookie = true
	}

	return userOutput, unsetTmpUserIDCookie, nil
}

func (s *service) RefreshToken(ctx context.Context, token string) (*erp.UserOutput, error) {
	claims := &auth.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})
	s.logger.Log("token", token)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, erp.ErrUnauthorized("%s", err)
		}
		return nil, erp.ErrBadRequest("%s", err)
	}
	if !tkn.Valid {
		return nil, erp.ErrUnauthorized("%s", err)
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return nil, erp.ErrBadRequest("%s", err)
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(auth.JwtKey)
	if err != nil {
		return nil, erp.ErrInternal("%s", err)
	}

	userOutput := &erp.UserOutput{}
	userOutput.Token = tokenString
	userOutput.ExpirationTime = expirationTime

	return userOutput, nil
}

func (s *service) GetUsers(ctx context.Context) ([]*erp.User, error) {
	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		return nil, erp.ErrInternal("%s", err)
	}
	return users, nil
}
