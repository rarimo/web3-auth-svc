package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/internal/service/requests"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewAuthorizeRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		address   = req.Data.ID
		signature = req.Data.Attributes.Signature

		log = Log(r).WithFields(map[string]any{
			"address":   address,
			"signature": signature,
		})
	)

	err = AuthVerifier(r).VerifySignature(signature, address)
	if err != nil {
		log.WithError(err).Info("Failed to verify signature")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	access, refresh, aexp, rexp, err := issueJWTs(r, address)
	if err != nil {
		log.WithError(err).Error("failed to issue JWTs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	Cookies(r).SetAccessToken(w, access, aexp)
	Cookies(r).SetRefreshToken(w, refresh, rexp)
	ape.Render(w, newTokenResponse(address, access, refresh))
}

func newTokenResponse(address, access, refresh string) resources.TokenResponse {
	return resources.TokenResponse{
		Data: resources.Token{
			Key: resources.Key{
				ID:   address,
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
}

func issueJWTs(r *http.Request, address string) (access, refresh string, aexp, rexp time.Time, err error) {
	access, aexp, err = JWT(r).IssueJWT(
		&jwt.AuthClaim{
			Address: address,
			Type:    jwt.AccessTokenType,
		},
	)
	if err != nil {
		return "", "", aexp, rexp, fmt.Errorf("failed to issue JWT access token: %w", err)
	}

	refresh, rexp, err = JWT(r).IssueJWT(
		&jwt.AuthClaim{
			Address: address,
			Type:    jwt.RefreshTokenType,
		},
	)
	if err != nil {
		return "", "", aexp, rexp, fmt.Errorf("failed to issue JWT access token: %w", err)
	}

	return access, refresh, aexp, rexp, nil
}
