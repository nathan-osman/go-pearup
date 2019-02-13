package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

var errUnsupportedProvider = errors.New("provider is not yet supported")

// FetchPicture attempts to retrieve a picture using the provided context.
func (a *Auth) FetchPicture(provider Provider, ctx context.Context, filename string) error {
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var url string
	switch provider {
	case ProviderFacebook:
		u, err := a.FacebookUserPicture(ctx)
		if err != nil {
			return err
		}
		url = u
	case ProviderGoogle:
		url = a.GooglePicture(ctx)
	default:
		return errUnsupportedProvider
	}
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = io.Copy(f, r.Body)
	return err
}
