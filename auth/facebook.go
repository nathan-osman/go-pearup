package auth

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/sling"
)

const facebookAPI = "https://graph.facebook.com/v3.0/"

// FacebookUserKey provides access to the Facebook user object in a request.
const FacebookUserKey = "facebook"

func (a *Auth) facebookLoginSucceeded(w http.ResponseWriter, r *http.Request) {
	fbUser, err := facebook.UserFromContext(r.Context())
	if err != nil {
		a.loginFailedHandler.ServeHTTP(w, withContext(r, ErrorKey, err))
		return
	}
	a.loginSucceededHandler.ServeHTTP(w, withContext(r, FacebookUserKey, fbUser))
}

// FacebookLoginHandler returns an HTTP handler for beginning Facebook auth.
func (a *Auth) FacebookLoginHandler() http.Handler {
	return facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.LoginHandler(
			a.oauth2Configs[ProviderFacebook],
			http.HandlerFunc(a.loginFailed),
		),
	)
}

// FacebookCallbackHandler returns an HTTP handler for completing Facebook auth.
func (a *Auth) FacebookCallbackHandler() http.Handler {
	return facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.CallbackHandler(
			a.oauth2Configs[ProviderFacebook],
			http.HandlerFunc(a.facebookLoginSucceeded),
			http.HandlerFunc(a.loginFailed),
		),
	)
}

// FacebookUserPicture retrieves the URL of a Facebook user's picture.
func (a *Auth) FacebookUserPicture(ctx context.Context) (string, error) {
	t, err := oauth2.TokenFromContext(ctx)
	if err != nil {
		return "", err
	}
	var (
		client  = a.oauth2Configs[ProviderFacebook].Client(ctx, t)
		picture = &struct {
			URL string `json:"url"`
		}{}
	)
	_, err = sling.
		New().
		Client(client).
		Base(facebookAPI).
		Set("Accept", "application/json").
		Get("me/picture").
		ReceiveSuccess(picture)
	return picture.URL, err
}
