package auth

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Auther interface {
	Auth() *Client
}

func NewAuther(getter kv.Getter) Auther {
	return &auther{
		getter: getter,
	}
}

type auther struct {
	once   comfig.Once
	getter kv.Getter
}

func (c *auther) Auth() *Client {
	return c.once.Do(func() interface{} {
		var cfg = struct {
			Addr    string `fig:"addr,required"`
			Enabled bool   `fig:"enabled"`
		}{
			Enabled: true,
		}

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "auth")).
			Please()

		if err != nil {
			panic(errors.WithMessage(err, "failed to figure out"))
		}

		return &Client{
			Client:  &http.Client{},
			Addr:    cfg.Addr,
			Enabled: cfg.Enabled,
		}
	}).(*Client)
}
