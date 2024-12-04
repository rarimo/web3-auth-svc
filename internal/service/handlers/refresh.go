package handlers

import (
	"net/http"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/resources"
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

	access, aexp, err := JWT(r).IssueJWT(
		&jwt.AuthClaim{
			Nullifier: claim.Nullifier,
			Type:      jwt.AccessTokenType,
		},
	)

	if err != nil {
		Log(r).WithError(err).WithField("user", claim.Nullifier).Error("failed to issuer JWT token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	refresh, rexp, err := JWT(r).IssueJWT(
		&jwt.AuthClaim{
			Nullifier: claim.Nullifier,
			Type:      jwt.RefreshTokenType,
		},
	)

	if err != nil {
		Log(r).WithError(err).WithField("user", claim.Nullifier).Error("failed to issuer JWT token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.TokenResponse{
		Data: resources.Token{
			Key: resources.Key{
				ID:   claim.Nullifier,
				Type: resources.TOKEN,
			},
			Attributes: resources.TokenAttributes{
				AccessToken: resources.Jwt{
					Token:     access,
					TokenType: string(jwt.AccessTokenType),
				},
				RefreshToken: resources.Jwt{
					Token:     refresh,
					TokenType: string(jwt.RefreshTokenType),
				},
			},
		},
	}

	Cookies(r).SetAccessToken(w, access, aexp)
	Cookies(r).SetRefreshToken(w, refresh, rexp)
	ape.Render(w, resp)
}
