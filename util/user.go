// Package spotify provides utilties for interfacing
// with Spotify's Web API.
package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

// User struct- used for making queries that require
// authentication access to a User's account
type User struct {
	http    *http.Client
	baseURL string

	autoRetry      bool
	acceptLanguage string

	auth   authenticator
	scopes scope
}

var (
	userExists = false
	ch         = make(chan *User)
)

// Get all scopes for a user
// See scopes.go for details on Scopes
func (u *User) GetScopes() []string {
	return u.scopes.Scopes
}

// Create and authenticate new User
// Defaults to using all (relevant) scopes if none are declared
func NewUser(client_key, secret_key string, scopes ...string) (*User, bool) {
	ok := true

	if len(scopes) == 0 {
		scopes = []string{
			ScopeUserModifyPlaybackState,
			ScopeUserReadPlaybackState,
			ScopeUserLibraryModify,
			ScopeUserLibraryRead,
			ScopeUserFollowModify,
			ScopePlaylistModifyPublic,
			ScopePlaylistModifyPrivate,
			ScopeUserFollowRead}
	}

	const redirectURI = "http://localhost:8080/callback"
	auth := newAuthenticator(
		"http://localhost:8080/callback", scopes...)

	state := ""

	// Define callback function
	completeAuth := func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			ok = false
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			// log.Fatalf("State mismatch: %s != %s\n", st, state)
			ok = false
		}
		// use the token to get an authenticated client
		client := auth.newClient(tok)

		fmt.Fprintf(w, "Login Completed!")
		ch <- &client
	}

	// First start an HTTP server
	// Only set HTTP callbacks for the first user
	if !userExists {
		http.HandleFunc("/callback", completeAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		})
	}
	go http.ListenAndServe(":8080", nil)

	auth.setAuthInfo(client_key, secret_key)

	// Use if you want login dialog to show every time
	// url := auth.authURLWithDialog(state)

	// Use if you want login dialog to show only when necessary
	url := auth.authURL(state)

	browser.OpenURL(url)

	// wait for auth to complete
	user := <-ch
	user.auth = auth
	user.scopes = scope{Scopes: scopes}

	userExists = true

	return user, ok
}

// Source for everything below: https://github.com/zmb3/spotify/

// errorStruct represents an error returned by the Spotify Web API.
type errorStruct struct {
	// A short description of the error.
	Message string `json:"message"`
	// The HTTP status code.
	Status int `json:"status"`
}

// return error message
func (e errorStruct) Error() string {
	return e.Message
}

// decodeError decodes an Error from an io.Reader.
func (u *User) decodeError(resp *http.Response) error {
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBody) == 0 {
		return fmt.Errorf("spotify: HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E errorStruct `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return fmt.Errorf("spotify: couldn't decode error: (%d) [%s]", len(responseBody), responseBody)
	}

	if e.E.Message == "" {
		// Some errors will result in there being a useful status-code but an
		// empty message, which will confuse the user (who only has access to
		// the message and not the code). An example of this is when we send
		// some of the arguments directly in the HTTP query and the URL ends-up
		// being too long.

		e.E.Message = fmt.Sprintf("spotify: unexpected HTTP %d: %s (empty error)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return e.E
}

// shouldRetry determines whether the status code indicates that the
// previous operation should be retried at a later time
func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

// isFailure determines whether the code indicates failure
func isFailure(code int, validCodes []int) bool {
	for _, item := range validCodes {
		if item == code {
			return false
		}
	}
	return true
}

// `execute` executes a non-GET request. `needsStatus` describes other HTTP
// status codes that will be treated as success. Note that we allow all 200s
// even if there are additional success codes that represent success.
func (u *User) execute(req *http.Request, result interface{}, needsStatus ...int) error {
	if u.acceptLanguage != "" {
		req.Header.Set("Accept-Language", u.acceptLanguage)
	}
	for {
		resp, err := u.http.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if u.autoRetry && shouldRetry(resp.StatusCode) {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if (resp.StatusCode >= 300 ||
			resp.StatusCode < 200) &&
			isFailure(resp.StatusCode, needsStatus) {
			return u.decodeError(resp)
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return err
			}
		}
		break
	}
	return nil
}

// get retry duration
func retryDuration(resp *http.Response) time.Duration {
	raw := resp.Header.Get("Retry-After")
	if raw == "" {
		return time.Second * 5
	}
	seconds, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return time.Second * 5
	}
	return time.Duration(seconds) * time.Second
}

// execute a get request
func (u *User) get(url string, result interface{}) error {
	for {
		req, err := http.NewRequest("GET", url, nil)
		if u.acceptLanguage != "" {
			req.Header.Set("Accept-Language", u.acceptLanguage)
		}
		if err != nil {
			return err
		}
		resp, err := u.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == 429 && u.autoRetry {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			return u.decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}

// token gets the client's current token.
func (u *User) token() (*oauth2.Token, error) {
	transport, ok := u.http.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("spotify: client not backed by oauth2 transport")
	}
	t, err := transport.Source.Token()
	if err != nil {
		return nil, err
	}
	return t, nil
}
