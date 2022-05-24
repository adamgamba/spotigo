package spotigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type accessToken struct {
	Token string `json:"access_token"`
}

// Get authentication token for all non-user account queries
func getToken(client string, secret string) (string, error) {
	var err error
	res, fetchErr := http.PostForm("https://accounts.spotify.com/api/token", url.Values{"grant_type": {"client_credentials"}, "client_id": {client}, "client_secret": {secret}})

	if fetchErr != nil {
		err = fetchErr
	}

	if res.StatusCode != 200 {
		err = errors.New("Status code: " + fmt.Sprint(res.StatusCode))
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		err = readErr
	}

	structure := accessToken{}
	jsonErr := json.Unmarshal(body, &structure)

	if jsonErr != nil {
		err = jsonErr
	}
	return structure.Token, err
}
