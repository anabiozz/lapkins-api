package auth

import (
	"net/http"
	"strings"

	"github.com/anabiozz/lapkins-api/pkg/cookies"
	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("secret")

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Subject string `json:"subject"`
	UserID  string `json:"user_id"`
	jwt.StandardClaims
}

func Check(token string) (*Claims, error) {
	token, err := stripBearerPrefixFromTokenString(token)
	if err != nil {
		return nil, err
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}
	return tkn.Claims.(*Claims), nil
}

func GetUserID(r *http.Request) (string, error) {
	token, err := stripBearerPrefixFromTokenString(r.Header.Get("Authorization"))
	if err != nil {
		return "", err
	}
	if token == "" {
		token, err = cookies.GetCookieValue(r, "token")
		if err != nil {
			if err == http.ErrNoCookie {
				tmpUserID, err := cookies.GetCookieValue(r, "tmp-user-id")
				if err != nil {
					if err != http.ErrNoCookie {
						return "", err
					}
				}
				return tmpUserID, nil
			} else {
				return "", err
			}
		}
	}
	if token != "" {
		claim, err := Check(token)
		if err != nil {
			return "", err
		}
		return claim.UserID, nil
	}
	return "", nil
}

func GetToken(r *http.Request) (string, error) {
	token, err := stripBearerPrefixFromTokenString(r.Header.Get("Authorization"))
	if err != nil {
		return "", err
	}
	if token == "" {
		token, err = cookies.GetCookieValue(r, "token")
		if err != nil {
			if err != http.ErrNoCookie {
				return "", err
			}
		}
	}
	return token, nil
}

func stripBearerPrefixFromTokenString(token string) (string, error) {
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:], nil
	}
	return token, nil
}
