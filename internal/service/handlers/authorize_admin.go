package handlers

import (
	"net/http"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func AuthorizeAdmin(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewAuthorizeAdminRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !Admins(r).Disabled && !Admins(r).VerifyAdmin(req.Data.Attributes.Password) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	access, aexp, err := JWT(r).IssueJWT(
		&jwt.AuthClaim{
			Address: "",
			Type:    jwt.AccessTokenType,
			IsAdmin: true,
		},
	)
	if err != nil {
		Log(r).WithError(err).Error("failed to issue JWT access token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	Cookies(r).SetAccessToken(w, access, aexp)
	ape.Render(w, newTokenResponse(req.Data.ID, access, ""))
}
