package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rarimo/web3-auth-svc/internal/cookies"
	"github.com/rarimo/web3-auth-svc/internal/jwt"
	"github.com/rarimo/web3-auth-svc/resources"
)

const (
	FullValidatePath = "integrations/web3-auth-svc/v1/validate"
)

type Client struct {
	*http.Client
	Addr    string
	Enabled bool
}

func (a *Client) ValidateJWT(r *http.Request) (claims []resources.Claim, err error) {
	validationURL, err := url.JoinPath(a.Addr, FullValidatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	req, err := http.NewRequest("GET", validationURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(jwt.AuthorizationHeaderName, r.Header.Get(jwt.AuthorizationHeaderName))
	req.Header.Set(cookies.CookieHeaderName, r.Header.Get(cookies.CookieHeaderName))

	resp, err := a.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute validate request: %w", err)
	}

	defer resp.Body.Close()

	body := resources.ValidationResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response body: %w", err)
	}

	return body.Data.Attributes.Claims, nil
}
