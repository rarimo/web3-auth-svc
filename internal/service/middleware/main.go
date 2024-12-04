package middleware

import (
	"net/http"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/internal/service/handlers"
	"github.com/rarimo/web3-auth-svc/pkg"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AuthMiddleware(issuer *jwt.JWTIssuer, log *logan.Entry, tokenType jwt.TokenType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := pkg.GetToken(r, tokenType)
			if err != nil {
				log.WithError(err).Debug("failed to get token")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			claim, err := issuer.ValidateJWT(token)
			if err != nil {
				log.WithError(err).Debug("failed validate bearer token")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			next.ServeHTTP(w, r.WithContext(handlers.CtxClaim(claim)(r.Context())))
		})
	}
}
