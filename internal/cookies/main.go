package cookies

import (
	"net/http"
	"time"

	"github.com/rarimo/web3-auth-svc/internal/jwt"
)

const (
	CookieHeaderName = "Cookie"
)

type Cookies struct {
	Domain   string
	Secure   bool
	SameSite int
}

func (c *Cookies) SetAccessToken(w http.ResponseWriter, token string, exp time.Time) {
	cookie := &http.Cookie{
		Name:     jwt.AccessTokenType.String(),
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: http.SameSite(c.SameSite),
		Domain:   c.Domain,
		Expires:  exp,
	}

	http.SetCookie(w, cookie)
}

func (c *Cookies) SetRefreshToken(w http.ResponseWriter, token string, exp time.Time) {
	cookie := &http.Cookie{
		Name:     jwt.RefreshTokenType.String(),
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: http.SameSite(c.SameSite),
		Domain:   c.Domain,
		Expires:  exp,
	}

	http.SetCookie(w, cookie)
}

func (c *Cookies) ClearTokensCookies(w http.ResponseWriter) {
	refreshCookie := &http.Cookie{
		Name:     jwt.RefreshTokenType.String(),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Secure,
		MaxAge:   -1,
		SameSite: http.SameSite(c.SameSite),
		Domain:   c.Domain,
	}

	http.SetCookie(w, refreshCookie)

	accessCookie := &http.Cookie{
		Name:     jwt.AccessTokenType.String(),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Secure,
		MaxAge:   -1,
		SameSite: http.SameSite(c.SameSite),
		Domain:   c.Domain,
	}

	http.SetCookie(w, accessCookie)
}
