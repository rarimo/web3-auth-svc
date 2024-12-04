package pkg

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func GetBearer(r *http.Request) (string, error) {
	authHeader := r.Header.Get(jwt.AuthorizationHeaderName)
	authHeaderSplit := strings.Split(authHeader, jwt.BearerTokenPrefix)

	if len(authHeaderSplit) != 2 {
		return "", ErrInvalidToken
	}

	return authHeaderSplit[1], nil
}

func SetBearer(r *http.Request, token string) {
	r.Header.Set(jwt.AuthorizationHeaderName, fmt.Sprintf("%s %s", jwt.BearerTokenPrefix, token))
}

func GetCookie(r *http.Request, typ jwt.TokenType) (string, error) {
	cookie, err := r.Cookie(typ.String())
	if err != nil {
		return "", ErrInvalidToken
	}

	return cookie.Value, nil
}

func GetToken(r *http.Request, typ jwt.TokenType) (string, error) {
	if token, err := GetBearer(r); err == nil {
		return token, nil
	}

	if token, err := GetCookie(r, typ); err == nil {
		return token, nil
	}

	return "", ErrInvalidToken
}
