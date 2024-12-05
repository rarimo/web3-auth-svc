package jwt

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gotest.tools/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	prv := make([]byte, 64)
	if _, err := rand.Read(prv); err != nil {
		panic(err)
	}

	fmt.Println(hexutil.Encode(prv))
}

func TestJWT(t *testing.T) {
	issuer := JWTIssuer{
		prv:               make([]byte, 64),
		accessExpiration:  time.Hour,
		refreshExpiration: time.Hour,
	}

	_, err := rand.Read(issuer.prv)
	assert.NilError(t, err)

	jwt, _, err := issuer.IssueJWT(
		&AuthClaim{
			Address: "0x31ba24c27a7d9b14fef5e48a26e79566525646ff",
			Type:    AccessTokenType,
		},
	)
	assert.NilError(t, err)

	claim, err := issuer.ValidateJWT(jwt)
	assert.NilError(t, err)

	assert.Equal(t, claim.Address, "0x31ba24c27a7d9b14fef5e48a26e79566525646ff")
	assert.Equal(t, claim.Type, AccessTokenType)
}
