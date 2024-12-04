package requests

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/web3-auth-svc/internal/zkp"
)

func GetPathNullifier(r *http.Request) (string, error) {
	nullifier := strings.ToLower(chi.URLParam(r, "nullifier"))
	return nullifier, validation.Errors{
		"nullifier": validation.Validate(nullifier, validation.Required, validation.Match(zkp.NullifierRegexp)),
	}.Filter()
}
