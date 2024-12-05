package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationHeaderName = "Authorization"
	BearerTokenPrefix       = "Bearer "
)

type JWTIssuer struct {
	prv               []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func (i *JWTIssuer) IssueJWT(claim *AuthClaim) (token string, exp time.Time, err error) {
	exp = time.Now().UTC()

	switch claim.Type {
	case AccessTokenType:
		exp = exp.Add(i.accessExpiration)
	case RefreshTokenType:
		exp = exp.Add(i.refreshExpiration)
	}

	raw := (&RawJWT{make(jwt.MapClaims)}).SetAddress(claim.Address).SetTokenType(claim.Type).SetExpirationTimestamp(exp)

	if claim.IsAdmin {
		raw.SetIsAdmin(true)
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, raw.claims).SignedString(i.prv)
	return
}

func (i *JWTIssuer) ValidateJWT(str string) (claim *AuthClaim, err error) {
	var token *jwt.Token

	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return i.prv, nil
	}

	if token, err = jwt.Parse(str, key, jwt.WithExpirationRequired()); err != nil {
		return nil, err
	}

	var (
		raw RawJWT
		ok  bool
	)
	if raw.claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, errors.New("failed to unwrap claims")
	}

	claim = &AuthClaim{}

	claim.Address, ok = raw.Address()
	if !ok {
		return nil, errors.New("invalid address: failed to parse")
	}

	claim.IsAdmin, _ = raw.IsAdmin()

	claim.Type, ok = raw.TokenType()
	if !ok {
		return nil, errors.New("invalid token type: failed to parse")
	}

	return
}
