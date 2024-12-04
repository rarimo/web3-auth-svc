package config

import (
	"github.com/rarimo/web3-auth-svc/internal/challenger"
	"github.com/rarimo/web3-auth-svc/internal/cookies"
	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	comfig.Logger
	comfig.Listenerer

	jwt.Jwter
	cookies.Cookier
	challenger.AuthVerifierer
}

type config struct {
	comfig.Logger
	comfig.Listenerer

	jwt.Jwter
	cookies.Cookier
	challenger.AuthVerifierer

	admin  comfig.Once
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:         getter,
		Listenerer:     comfig.NewListenerer(getter),
		Logger:         comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Jwter:          jwt.NewJwter(getter),
		Cookier:        cookies.NewCookier(getter),
		AuthVerifierer: challenger.NewAuthVerifierer(getter),
	}
}
