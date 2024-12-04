package handlers

import (
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rarimo/web3-auth-svc/internal/service/requests"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RequestChallenge(w http.ResponseWriter, r *http.Request) {
	nullifier, err := requests.GetPathNullifier(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifierBytes, err := hexutil.Decode(nullifier)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	challenge, err := Verifier(r).Challenge(new(big.Int).SetBytes(nullifierBytes).String())
	if err != nil {
		Log(r).WithError(err).Error("failed to generate challenge")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.ChallengeResponse{
		Data: resources.Challenge{
			Key: resources.Key{
				ID:   nullifier,
				Type: resources.CHALLENGE,
			},
			Attributes: resources.ChallengeAttributes{
				Challenge: challenge,
			},
		},
	}

	ape.Render(w, resp)
}
