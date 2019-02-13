package auth

import (
	"net/http"
)

type Provider string

const (
	ProviderFacebook  Provider = "facebook"
	ProviderGoogle             = "google"
	ProviderMicrosoft          = "microsoft"
)

// ProviderConfig stores configuration information for a specific provider.
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Config stores configuration information for authentication.
type Config struct {
	ProviderConfigs       map[Provider]*ProviderConfig
	LoginSucceededHandler http.Handler
	LoginFailedHandler    http.Handler
}
