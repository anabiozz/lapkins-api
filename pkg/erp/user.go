package erp

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errMustProvidePhoneOrEmail = errors.New("must provide phone or email")
	errInvalidGender           = errors.New("invalid Gender")
)

type User struct {
	ID               primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName        string             `bson:"first_name" json:"first_name"`
	Patronymic       string             `bson:"patronymic" json:"patronymic"`
	LastName         string             `bson:"last_name" json:"last_name"`
	Gender           Gender             `bson:"gender" json:"gender"`
	Password         string             `bson:"password" json:"password"`
	Login            string             `bson:"login" json:"login"`
	Email            string             `bson:"email" json:"email"`
	IsActive         bool               `bson:"is_active" json:"is_active"`
	RegistrationDate time.Time          `bson:"registration_date" json:"registration_date"`
	Phone            int64              `bson:"phone" json:"phone"`
	Birthday         time.Time          `bson:"birthday" json:"birthday"`
	ConstDiscount    uint8              `bson:"const_discount" json:"const_discount"`
}

type UserInput struct {
	Password   string    `json:"password"`
	Login      string    `json:"login"`
	FirstName  string    `json:"first_name"`
	Patronymic string    `json:"patronymic"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      int64     `json:"phone"`
	Birthday   time.Time `json:"birthday"`
	Gender     Gender    `json:"gender"`
}

type UserOutput struct {
	ID             string    `json:"id"`
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type Response struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func (in *UserInput) Init(createdTime time.Time) *User {
	createdTime = createdTime.Round(time.Second)
	u := &User{
		Login:            strings.TrimSpace(in.Login),
		FirstName:        strings.TrimSpace(in.FirstName),
		Patronymic:       strings.TrimSpace(in.Patronymic),
		LastName:         strings.TrimSpace(in.LastName),
		Email:            strings.TrimSpace(in.Email),
		IsActive:         true,
		Password:         in.Password,
		RegistrationDate: createdTime,
		Phone:            in.Phone,
		Birthday:         in.Birthday.Round(time.Second),
		Gender:           in.Gender,
		ConstDiscount:    0,
	}
	return u
}

func (in *UserInput) Validate() error {
	if strings.TrimSpace(in.Email) == "" && in.Phone <= 0 {
		return errMustProvidePhoneOrEmail
	}
	if !in.Gender.IsValid() {
		return errInvalidGender
	}
	return nil
}

type Gender string

func (g Gender) MarshalJSON() ([]byte, error) {
	if g == "" {
		return []byte("null"), nil
	}

	return json.Marshal(string(g))
}

const (
	GenderNotSpecified Gender = ""
	GenderMale         Gender = "M"
	GenderFemale       Gender = "F"
)

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale:
	case GenderFemale:
	case GenderNotSpecified:
	default:
		return false
	}
	return true
}
