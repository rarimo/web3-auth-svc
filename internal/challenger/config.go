package challenger

import (
	"fmt"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type AuthVerifierer interface {
	AuthVerifier() *AuthVerifier
}

func NewAuthVerifierer(getter kv.Getter) AuthVerifierer {
	return &authVerifierer{
		getter: getter,
	}
}

type authVerifierer struct {
	once   comfig.Once
	getter kv.Getter
}

func (v *authVerifierer) AuthVerifier() *AuthVerifier {
	return v.once.Do(func() interface{} {
		cfg := struct {
			Disabled bool `fig:"disabled"`
		}{}

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(v.getter, "auth_verifier")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out: %w", err))
		}

		return &AuthVerifier{
			Disabled:   cfg.Disabled,
			challenges: make(map[string]*Challenge),
		}
	}).(*AuthVerifier)
}
