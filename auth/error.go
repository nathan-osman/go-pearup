package auth

import (
	"net/http"

	"github.com/dghubble/gologin"
)

// ErrorKey provides access to an errror.
const ErrorKey = "error"

func (a *Auth) loginFailed(w http.ResponseWriter, r *http.Request) {
	a.loginFailedHandler.ServeHTTP(
		w, withContext(r, ErrorKey, gologin.ErrorFromContext(r.Context())),
	)
}
