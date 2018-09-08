package auth

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
)

type key int

const (
	// ErrorKey provides access to an errror.
	ErrorKey key = iota

	// FacebookUserKey provides access to the Facebook user struct.
	FacebookUserKey
)

func (a *Auth) loginSucceeded(w http.ResponseWriter, r *http.Request) {
	fbUser, err := facebook.UserFromContext(r.Context())
	if err != nil {
		r = r.WithContext(context.WithValue(r.Context(), ErrorKey, err))
		a.loginFailedHandler.ServeHTTP(w, r)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), FacebookUserKey, fbUser))
	a.loginSucceededHandler.ServeHTTP(w, r)
}

func (a *Auth) loginFailed(w http.ResponseWriter, r *http.Request) {
	err := gologin.ErrorFromContext(r.Context()).Error()
	r = r.WithContext(context.WithValue(r.Context(), ErrorKey, err))
	a.loginFailedHandler.ServeHTTP(w, r)
}
