package auth

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/google"
	"google.golang.org/api/oauth2/v2"
)

// GoogleUserKey provides access to the Google user object in a request.
const GoogleUserKey = "google"

func (a *Auth) googleLoginSucceeded(w http.ResponseWriter, r *http.Request) {
	gUser, err := google.UserFromContext(r.Context())
	if err != nil {
		a.loginFailedHandler.ServeHTTP(w, withContext(r, ErrorKey, err))
		return
	}
	a.loginSucceededHandler.ServeHTTP(w, withContext(r, GoogleUserKey, gUser))
}

// GoogleLoginHandler returns an HTTP handler for beginning Google auth.
func (a *Auth) GoogleLoginHandler() http.Handler {
	return google.StateHandler(
		gologin.DebugOnlyCookieConfig,
		google.LoginHandler(
			a.oauth2Configs[ProviderGoogle],
			http.HandlerFunc(a.loginFailed),
		),
	)
}

// GoogleCallbackHandler returns an HTTP handler for completing Google auth.
func (a *Auth) GoogleCallbackHandler() http.Handler {
	return google.StateHandler(
		gologin.DebugOnlyCookieConfig,
		google.CallbackHandler(
			a.oauth2Configs[ProviderGoogle],
			http.HandlerFunc(a.googleLoginSucceeded),
			http.HandlerFunc(a.loginFailed),
		),
	)
}

// GooglePicture retrieves the URL of the user's picture.
func (a *Auth) GooglePicture(ctx context.Context) string {
	return ctx.Value(GoogleUserKey).(*oauth2.Userinfoplus).Picture
}
