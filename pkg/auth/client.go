package auth

import (
	"encoding/json"
	"fmt"
	"github.com/rarimo/web3-auth-svc/internal/cookies"
	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"net/http"

	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	FullValidatePath = "integrations/web3-auth-svc/v1/validate"
)

type Client struct {
	*http.Client
	Addr string
}

func (a *Client) ValidateJWT(r *http.Request) (claims []resources.Claim, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.Addr, FullValidatePath), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set(jwt.AuthorizationHeaderName, r.Header.Get(jwt.AuthorizationHeaderName))
	req.Header.Set(cookies.CookieHeaderName, r.Header.Get(cookies.CookieHeaderName))

	resp, err := a.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute validate request")
	}

	defer resp.Body.Close()

	body := resources.ValidationResultResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall response body")
	}

	return body.Data.Attributes.Claims, nil
}
