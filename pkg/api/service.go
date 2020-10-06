package api

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

type Storage interface {
	GetProduct(ctx context.Context, sku string) (*model.Product, error)
	GetProducts(ctx context.Context) ([]*model.Product, error)
	GetCategories(ctx context.Context) ([]*model.Category, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string, category *model.Category) error
	RemoveCategory(ctx context.Context, sku string, category *model.Category) error
	UpdateProduct(ctx context.Context, product *model.Product) error

	AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error)
	IncreaseProductQuantity(ctx context.Context, userID string, sku string) error
	DecreaseProductQuantity(ctx context.Context, userID string, sku string) error
	RemoveProduct(ctx context.Context, userID string, sku string) error
	LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error)
	AddOrder(ctx context.Context, order *model.Order) error

	RegisterUser(ctx context.Context, user *model.User) (string, error)
	Login(ctx context.Context, email string, phone int64, tmpUserID string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
}

type Service interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetProducts(ctx context.Context) ([]*model.Product, error)
	GetProduct(ctx context.Context, sku string) (*model.Product, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string, category *model.Category) error
	RemoveCategory(ctx context.Context, sku string, category *model.Category) error
	UpdateProduct(ctx context.Context, product *model.Product) error

	AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error)
	DecreaseProductQuantity(ctx context.Context, userID string, sku string) error
	IncreaseProductQuantity(ctx context.Context, userID string, sku string) error
	LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error)
	RemoveProduct(ctx context.Context, userID string, sku string) error

	AddOrder(ctx context.Context, order *model.Order) error

	Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error)
	Login(ctx context.Context, input *model.UserInput, tmpUserID string) (*model.UserOutput, bool, error)
	RefreshToken(ctx context.Context, token string) (*model.UserOutput, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
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

func (s *service) GetCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.storage.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	product, err := s.storage.GetProduct(ctx, sku)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *service) GetProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.storage.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *service) UpdateProduct(ctx context.Context, product *model.Product) error {
	err := s.storage.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
	err := s.storage.AddAttribute(ctx, sku, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	err := s.storage.RemoveAttribute(ctx, sku, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	err := s.storage.AddCategory(ctx, sku, category)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	err := s.storage.RemoveCategory(ctx, sku, category)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error) {
	setTmpUserIDCookie, userID, err := s.storage.AddProductToCard(ctx, sku, userID, isLoggedIn, isTmpUserIDSet)
	if err != nil {
		return false, "", errBadRequest("%s", err)
	}
	return setTmpUserIDCookie, userID, nil
}

func (s *service) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	err := s.storage.DecreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		return errBadRequest("%s", err)
	}
	return nil
}

func (s *service) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	err := s.storage.IncreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		return errBadRequest("%s", err)
	}
	return nil
}

func (s *service) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	if userID == "" {
		return nil, errBadRequest("%s", "provided user id is empty")
	}
	cart, err := s.storage.LoadCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return cart, nil

}

func (s *service) RemoveProduct(ctx context.Context, userID string, sku string) error {
	if userID == "" {
		return errBadRequest("%s", "provided user id is empty")
	}
	err := s.storage.RemoveProduct(ctx, userID, sku)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddOrder(ctx context.Context, order *model.Order) error {
	err := s.storage.AddOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}


func (s *service) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
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

func (s *service) Login(ctx context.Context, input *model.UserInput, tmpUserID string) (*model.UserOutput, bool, error) {
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

func (s *service) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
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

func (s *service) GetUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		return nil, errInternal("%s", err)
	}
	return users, nil
}