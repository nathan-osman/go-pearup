package auth

import (
	"net/http"
)

// Config stores configuration information for authentication.
type Config struct {
	FacebookClientID              string
	FacebookClientSecret          string
	FacebookRedirectURL           string
	FacebookLoginSucceededHandler http.Handler
	FacebookLoginFailedHandler    http.Handler
}
