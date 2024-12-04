package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RequestChallenge(w http.ResponseWriter, r *http.Request) {
	address := strings.ToLower(chi.URLParam(r, "address"))

	challenge, err := AuthVerifier(r).Challenge(address)
	if err != nil {
		Log(r).WithError(err).Error("failed to generate challenge")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.ChallengeResponse{
		Data: resources.Challenge{
			Key: resources.Key{
				ID:   address,
				Type: resources.CHALLENGE,
			},
			Attributes: resources.ChallengeAttributes{
				Challenge: challenge,
			},
		},
	})
}
