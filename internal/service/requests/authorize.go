package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/web3-auth-svc/internal/zkp"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewAuthorizeRequest(r *http.Request) (*resources.AuthorizeRequest, error) {
	req := &resources.AuthorizeRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.New("failed to unmarshall request")
	}

	req.Data.ID = strings.ToLower(req.Data.ID)
	return req, validation.Errors{
		"data/id":   validation.Validate(req.Data.ID, validation.Required, validation.Match(zkp.NullifierRegexp)),
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.AUTHORIZE)),
	}.Filter()
}
