package auth

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/sling"
)

const (
	facebookAPI = "https://graph.facebook.com/v3.0/"
)

// UpdatePicture updates the picture for the user represented by the context, storing it in the specified file.
func (a *Auth) UpdatePicture(ctx context.Context, filename string) error {
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	t, err := oauth2.TokenFromContext(ctx)
	if err != nil {
		return err
	}
	r, err := sling.New().Base(facebookAPI).Get("me/picture").Request()
	if err != nil {
		return err
	}
	resp, err := a.oauth2Config.Client(ctx, t).Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	return nil
}
