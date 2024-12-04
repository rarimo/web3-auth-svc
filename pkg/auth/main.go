package auth

import "github.com/rarimo/web3-auth-svc/resources"

type Grant func(claim resources.Claim) bool

func Authenticates(claims []resources.Claim, grants ...Grant) bool {
	for _, claim := range claims {
		for _, grant := range grants {
			if grant(claim) {
				return true
			}
		}
	}

	return false
}
