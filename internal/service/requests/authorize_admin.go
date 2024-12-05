package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewAuthorizeAdminRequest(r *http.Request) (req resources.AuthorizeAdminRequest, err error) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("failed to unmarshall request")
	}

	return req, validation.Errors{
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.AUTHORIZE)),
	}.Filter()
}
