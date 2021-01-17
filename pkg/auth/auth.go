package auth

import (
	"net/http"
	"strings"

	"github.com/anabiozz/core/lapkins/pkg/cookies"
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
	token = stripBearerPrefixFromTokenString(token)
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

func GetUserID(r *http.Request) (string, bool, error) {
	token, err := GetToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			userID, err := cookies.GetCookieValue(r, "tmp-user-id")
			if err != nil {
				if err == http.ErrNoCookie {
					return "", false, nil
				}
				return "", false, err
			}
			return userID, false, nil
		} else {
			return "", false, err
		}
	}
	if token != "" {
		claim, err := Check(token)
		if err != nil {
			return "", false, err
		}
		return claim.UserID, true, nil
	}
	return "", false, nil
}

func GetToken(r *http.Request) (string, error) {
	token := stripBearerPrefixFromTokenString(r.Header.Get("Authorization"))
	if token == "" {
		token, err := cookies.GetCookieValue(r, "token")
		if err != nil {
			return "", err
		}
		return token, nil
	}
	return token, nil
}

func stripBearerPrefixFromTokenString(token string) string {
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:]
	}
	return token
}
