package handlers

import (
	"net/http"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	claim := Claim(r)
	if claim == nil {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	if claim.Type != jwt.RefreshTokenType {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	access, refresh, aexp, rexp, err := issueJWTs(r, claim.Address)
	if err != nil {
		Log(r).WithError(err).WithField("address", claim.Address).Error("failed to issue JWTs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	Cookies(r).SetAccessToken(w, access, aexp)
	Cookies(r).SetRefreshToken(w, refresh, rexp)
	ape.Render(w, newTokenResponse(claim.Address, access, refresh))
}
