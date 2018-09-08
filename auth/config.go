package auth

import (
	"net/http"
)

// Config stores configuration information for Facebook authentication.
type Config struct {
	ClientID              string
	ClientSecret          string
	RedirectURL           string
	LoginSucceededHandler http.Handler
	LoginFailedHandler    http.Handler
}
