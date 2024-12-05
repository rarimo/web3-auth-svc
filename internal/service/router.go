package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/rarimo/web3-auth-svc/internal/config"
	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/internal/service/handlers"
	"github.com/rarimo/web3-auth-svc/internal/service/middleware"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, cfg config.Config) {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxJWT(cfg.JWT()),
			handlers.CtxCookies(cfg.Cookies()),
			handlers.CtxAuthVerifier(cfg.AuthVerifier()),
			handlers.CtxAdmins(cfg.Admin()),
		),
	)

	r.Route("/integrations/web3-auth-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/authorize", func(r chi.Router) {
				r.Post("/admin", handlers.AuthorizeAdmin)
				r.Post("/", handlers.Authorize)
				r.Get("/{address}/challenge", handlers.RequestChallenge)
			})
			r.With(middleware.AuthMiddleware(cfg.Log(), jwt.AccessTokenType)).Get("/validate", handlers.Validate)
			r.With(middleware.AuthMiddleware(cfg.Log(), jwt.RefreshTokenType)).Get("/refresh", handlers.Refresh)
		})
	})

	cfg.Log().Info("Service started")
	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
