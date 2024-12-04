package auth

import (
	"github.com/rarimo/web3-auth-svc/resources"
)

func UserGrant(nullifier string) Grant {
	return func(claim resources.Claim) bool {
		return claim.Nullifier == nullifier
	}
}
