package util

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// Source: https://github.com/zmb3/spotify/

const (
	// AuthURL is the URL to Spotify Accounts Service's OAuth2 endpoint.
	authURL = "https://accounts.spotify.com/authorize"
	// TokenURL is the URL to the Spotify Accounts Service's OAuth2
	// token endpoint.
	tokenURL = "https://accounts.spotify.com/api/token"
)

// authenticator provides convenience functions for implementing the OAuth2 flow.
type authenticator struct {
	config  *oauth2.Config
	context context.Context
}

// NewAuthenticator creates an authenticator which is used to implement the
// OAuth2 authorization flow.  The redirectURL must exactly match one of the
// URLs specified in your Spotify developer account.
//
func newAuthenticator(redirectURL string, scopes ...string) authenticator {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	// disable HTTP/2 for DefaultClient, see: https://github.com/zmb3/spotify/issues/20
	tr := &http.Transport{
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr})
	return authenticator{
		config:  cfg,
		context: ctx,
	}
}

// SetAuthInfo overwrites the client ID and secret key used by the authenticator.
// You can use this if you don't want to store this information in environment variables.
func (a *authenticator) setAuthInfo(clientID, secretKey string) {
	a.config.ClientID = clientID
	a.config.ClientSecret = secretKey
}

// AuthURL returns a URL to the the Spotify Accounts Service's OAuth2 endpoint.
//
// State is a token to protect the user from CSRF attacks.  You should pass the
// same state to `Token`, where it will be validated.  For more info, refer to
// http://tools.ietf.org/html/rfc6749#section-10.12.
func (a authenticator) authURL(state string) string {
	return a.config.AuthCodeURL(state)
}

// AuthURLWithDialog returns the same URL as AuthURL, but sets show_dialog to true
func (a authenticator) authURLWithDialog(state string) string {
	return a.config.AuthCodeURL(state, oauth2.SetAuthURLParam("show_dialog", "true"))
}

// Token pulls an authorization code from an HTTP request and attempts to exchange
// it for an access token.  The standard use case is to call Token from the handler
// that handles requests to your application's redirect URL.
func (a authenticator) token(state string, r *http.Request) (*oauth2.Token, error) {
	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		return nil, errors.New("spotify: auth failed - " + e)
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New("spotify: didn't get access code")
	}
	actualState := values.Get("state")
	if actualState != state {
		return nil, errors.New("spotify: redirect state parameter doesn't match")
	}
	return a.config.Exchange(a.context, code)
}

// NewClient creates a Client that will use the specified access token for its API requests.
func (a authenticator) newClient(token *oauth2.Token) User {
	client := a.config.Client(a.context, token)
	return User{
		http:    client,
		baseURL: "https://api.spotify.com/v1/",
	}
}
