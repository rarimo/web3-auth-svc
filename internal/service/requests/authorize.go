package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/web3-auth-svc/internal/challenger"
	"github.com/rarimo/web3-auth-svc/resources"
)

func NewAuthorizeRequest(r *http.Request) (req resources.AuthorizeRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	req.Data.ID = strings.ToLower(req.Data.ID)
	return req, validation.Errors{
		"data/id":   validation.Validate(req.Data.ID, validation.Required, validation.Match(challenger.AddressRegexp)),
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.AUTHORIZE)),
		"data/attributes/signature": validation.Validate(
			req.Data.Attributes.Signature,
			validation.Required,
			validation.Match(challenger.SignatureRegexp)),
	}.Filter()
}

func newDecodeError(what string, err error) error {
	return validation.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}
