package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ClaimAddress             = "sub"
	ClaimTokenType           = "type"
	ClaimExpirationTimestamp = "exp"
)

type TokenType string

func (t TokenType) String() string {
	return string(t)
}

var (
	AccessTokenType  TokenType = "access"
	RefreshTokenType TokenType = "refresh"
)

// AuthClaim is a helper structure to organize all claims in one entity
type AuthClaim struct {
	Address string
	Type    TokenType
}

// RawJWT represents helper structure to provide setter and getter methods to work with JWT claims
type RawJWT struct {
	claims jwt.MapClaims
}

// Setters

func (r *RawJWT) SetAddress(sub string) *RawJWT {
	r.claims[ClaimAddress] = sub
	return r
}

func (r *RawJWT) SetExpirationTimestamp(expiration time.Time) *RawJWT {
	r.claims[ClaimExpirationTimestamp] = jwt.NewNumericDate(expiration)
	return r
}

func (r *RawJWT) SetTokenType(typ TokenType) *RawJWT {
	r.claims[ClaimTokenType] = typ
	return r
}

// Getters

func (r *RawJWT) Address() (res string, ok bool) {
	var val interface{}

	if val, ok = r.claims[ClaimAddress]; !ok {
		return
	}

	res, ok = val.(string)
	return
}

func (r *RawJWT) TokenType() (typ TokenType, ok bool) {
	var (
		val interface{}
		str string
	)

	if val, ok = r.claims[ClaimTokenType]; !ok {
		return
	}

	if str, ok = val.(string); !ok {
		return
	}

	return TokenType(str), true
}
