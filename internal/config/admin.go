package config

import (
	"crypto/sha256"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/web3-auth-svc/internal/challenger"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type Admin struct {
	Admin    string `fig:"password_hash,required"`
	Disabled bool   `fig:"disabled"`
}

func (c *config) Admin() *Admin {
	return c.admin.Do(func() interface{} {
		var cfg Admin

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "admin")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out admin config: %w", err))
		}

		if !cfg.Disabled {
			err := validation.Errors{
				"admin/password_hash": validation.Validate(cfg.Admin, validation.Required, validation.Match(challenger.SHA256HexRegexp)),
			}.Filter()
			if err != nil {
				panic(err)
			}
		}

		return &cfg
	}).(*Admin)
}

func (p *Admin) VerifyAdmin(pass string) bool {
	passHash := sha256.Sum256([]byte(pass))
	return slices.Equal(hexutil.MustDecode(p.Admin), passHash[:])
}
