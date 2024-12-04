package config

import (
	"github.com/rarimo/web3-auth-svc/internal/cookies"
	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/internal/zkp"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	comfig.Logger
	comfig.Listenerer
	jwt.Jwter
	zkp.Verifierer
	cookies.Cookier
}

type config struct {
	comfig.Logger
	comfig.Listenerer
	jwt.Jwter
	zkp.Verifierer
	cookies.Cookier
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Jwter:      jwt.NewJwter(getter),
		Verifierer: zkp.NewVerifierer(getter),
		Cookier:    cookies.NewCookier(getter),
	}
}
