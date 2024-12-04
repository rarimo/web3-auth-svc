package auth

import (
	"github.com/rarimo/web3-auth-svc/resources"
)

func UserGrant(address string) Grant {
	return func(claim resources.Claim) bool {
		return claim.Address == address
	}
}
