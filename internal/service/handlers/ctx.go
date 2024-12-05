package handlers

import (
	"context"
	"net/http"

	"github.com/rarimo/web3-auth-svc/internal/challenger"
	"github.com/rarimo/web3-auth-svc/internal/config"
	"github.com/rarimo/web3-auth-svc/internal/cookies"
	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	claimCtxKey
	jwtCtxKey
	cookiesCtxKey
	authVerifierCtxKey
	adminsCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxClaim(claim *jwt.AuthClaim) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, claimCtxKey, claim)
	}
}

func Claim(r *http.Request) *jwt.AuthClaim {
	return r.Context().Value(claimCtxKey).(*jwt.AuthClaim)
}

func CtxJWT(issuer *jwt.JWTIssuer) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, jwtCtxKey, issuer)
	}
}

func JWT(r *http.Request) *jwt.JWTIssuer {
	return r.Context().Value(jwtCtxKey).(*jwt.JWTIssuer)
}

func CtxCookies(cookies *cookies.Cookies) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, cookiesCtxKey, cookies)
	}
}

func Cookies(r *http.Request) *cookies.Cookies {
	return r.Context().Value(cookiesCtxKey).(*cookies.Cookies)
}

func CtxAuthVerifier(verifier *challenger.AuthVerifier) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, authVerifierCtxKey, verifier)
	}
}

func AuthVerifier(r *http.Request) *challenger.AuthVerifier {
	return r.Context().Value(authVerifierCtxKey).(*challenger.AuthVerifier)
}

func CtxAdmins(admins *config.Admin) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, adminsCtxKey, admins)
	}
}

func Admins(r *http.Request) *config.Admin {
	return r.Context().Value(adminsCtxKey).(*config.Admin)
}
