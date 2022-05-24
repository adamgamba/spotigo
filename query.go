package helloworld

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Query struct- stores developer's client and secret IDs
// and their API access token
type Query struct {
	client string
	secret string
	token  string `json:"access_token"`
}

// Constructor- create a new Query
func NewQuery(client string, secret string) (Query, error) {
	q := Query{client: client, secret: secret}

	token, authErr := getToken(client, secret)
	q.token = token

	return q, authErr
}

// Query Methods

// Get Album by URI
func (q Query) GetAlbumByURI(uri string) (Album, bool) {
	bytes, fetchErr := fetch(q.token, uri, "albums")
	album := Album{}
	ok := true
	if fetchErr != nil {
		ok = false
		fmt.Println("Fetch Error:", fetchErr)
	}

	jsonErr := json.Unmarshal(bytes, &album)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	return album, ok
}

// Get Artist by URI
func (q Query) GetArtistByURI(uri string) (Artist, bool) {
	bytes, fetchErr := fetch(q.token, uri, "artists")
	artist := Artist{}
	ok := true
	if fetchErr != nil {
		ok = false
		fmt.Println("Fetch Error:", fetchErr)
	}

	jsonErr := json.Unmarshal(bytes, &artist)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	return artist, ok
}

// Get Track by URI
func (q Query) GetTrackByURI(uri string) (Track, bool) {
	bytes, fetchErr := fetch(q.token, uri, "tracks")
	track := Track{}
	ok := true
	if fetchErr != nil {
		ok = false
		fmt.Println("Fetch Error:", fetchErr)
	}

	jsonErr := json.Unmarshal(bytes, &track)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	return track, ok
}

// Get Playlist by URI
func (q Query) GetPlaylistByURI(uri string) (Playlist, bool) {
	bytes, fetchErr := fetch(q.token, uri, "playlists")
	playlist := Playlist{}
	ok := true
	if fetchErr != nil {
		ok = false
		fmt.Println("Fetch Error:", fetchErr)
	}

	jsonErr := json.Unmarshal(bytes, &playlist)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	return playlist, ok
}

// General search function- searches Spotify platform and returns
// the first result that corresponds to returnType
// Possible returnTypes: "album","artist","playlist","track"
// The input parameter is a string search query; for the user to reliably get the
// track they're looking for, they need to include as much information in this
// string as possible. Misspelled or incomplete inputs will return a result,
// but the more information included, the more likely the result will be as intended
// Bad input example: "disco"
// Good input example: "Disco Man Remi Wolf"
func (q Query) search(input string, returnType string) ([]byte, error) {
	var err error
	reqClient := &http.Client{}
	baseURL := "https://api.spotify.com/v1/"

	input = url.QueryEscape(input)

	req, httpErr := http.NewRequest("GET", (baseURL + "search" + "?q=" + input + "&type=" + returnType), nil)
	req.Header.Add("Authorization", ("Bearer " + q.token))
	req.Header.Add("Content-Type", "application/json")

	if httpErr != nil {
		err = httpErr
	}

	res, fetchErr := reqClient.Do(req)
	if fetchErr != nil {
		err = fetchErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		err = httpErr
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = errors.New("Status code: " + fmt.Sprint(res.StatusCode))
	}
	return body, err
}

// Get Artist by name
// input is a string search query as described in the search function
func (q Query) GetArtistByName(input string) (Artist, bool) {
	bytes, searchErr := q.search(input, "artist")

	searchResult := searchResult{}
	ok := true
	if searchErr != nil {
		ok = false
		fmt.Println("Search Error:", searchErr)
	}

	jsonErr := json.Unmarshal(bytes, &searchResult)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	if searchResult.Artists.Total > 0 {
		artist := searchResult.Artists.Items[0]
		return artist, ok
	}

	// No Results Found
	return Artist{}, false
}

// Get Album by name
// input is a string search query as described in the search function
func (q Query) GetAlbumByName(input string) (Album, bool) {
	bytes, searchErr := q.search(input, "album")

	searchResult := searchResult{}
	ok := true
	if searchErr != nil {
		ok = false
		fmt.Println("Search Error:", searchErr)
	}

	jsonErr := json.Unmarshal(bytes, &searchResult)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	if searchResult.Albums.Total > 0 {
		album := searchResult.Albums.Items[0]
		return album, ok
	}

	// No Results Found
	return Album{}, false
}

// Get Track by name
// input is a string search query as described in the search function
func (q Query) GetTrackByName(input string) (Track, bool) {
	bytes, searchErr := q.search(input, "track")

	searchResult := searchResult{}
	ok := true

	if searchErr != nil {
		ok = false
		fmt.Println("Search Error:", searchErr)
	}

	jsonErr := json.Unmarshal(bytes, &searchResult)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	if searchResult.Tracks.Total > 0 {
		track := searchResult.Tracks.Items[0]
		return track, ok
	}

	// No Results Found
	return Track{}, false
}

// Get Playlist by name
// input is a string search query as described in the search function
func (q Query) GetPlaylistByName(input string) (Playlist, bool) {
	bytes, searchErr := q.search(input, "playlist")

	searchResult := searchResult{}
	ok := true

	if searchErr != nil {
		ok = false
		fmt.Println("Search Error:", searchErr)
	}

	jsonErr := json.Unmarshal(bytes, &searchResult)

	if jsonErr != nil {
		ok = false
		fmt.Println("JSON Error:", jsonErr)
	}

	if searchResult.Playlists.Total > 0 {
		playlist := searchResult.Playlists.Items[0]
		return playlist, ok
	}

	// No Results Found
	return Playlist{}, false
}

// Execute HTTP request
func fetch(token string, uri string, endpoint string) ([]byte, error) {
	var err error
	reqClient := &http.Client{}
	baseURL := "https://api.spotify.com/v1/"

	req, httpErr := http.NewRequest("GET", (baseURL + endpoint + "/" + uri), nil)
	req.Header.Add("Authorization", ("Bearer " + token))
	req.Header.Add("Content-Type", "application/json")
	if httpErr != nil {
		err = httpErr
	}

	res, fetchErr := reqClient.Do(req)
	if fetchErr != nil {
		err = fetchErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		err = httpErr
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = errors.New("Status code: " + fmt.Sprint(res.StatusCode))
	}

	return body, err
}

// Get multiple Tracks by URIs
func (q Query) GetTracksByURIs(uris ...string) ([]Track, bool) {
	tracks := make([]Track, 0)
	ok := true

	for _, x := range uris {
		track, thisOk := q.GetTrackByURI(x)
		ok = ok && thisOk
		tracks = append(tracks, track)
	}
	return tracks, ok
}

// Get multiple Tracks by names
func (q Query) GetTracksByNames(names ...string) ([]Track, bool) {
	tracks := make([]Track, 0)
	ok := true

	for _, x := range names {
		track, thisOk := q.GetTrackByName(x)
		ok = ok && thisOk
		tracks = append(tracks, track)
	}
	return tracks, ok
}

// Get multiple Albums by URIs
func (q Query) GetAlbumsByURIs(uris ...string) ([]Album, bool) {
	albums := make([]Album, 0)
	ok := true

	for _, x := range uris {
		album, thisOk := q.GetAlbumByURI(x)
		ok = ok && thisOk
		albums = append(albums, album)
	}
	return albums, ok
}

// Get multiple Albums by names
func (q Query) GetAlbumsByNames(names ...string) ([]Album, bool) {
	albums := make([]Album, 0)
	ok := true

	for _, x := range names {
		album, thisOk := q.GetAlbumByName(x)
		ok = ok && thisOk
		albums = append(albums, album)
	}
	return albums, ok
}

// Get multiple Artists by URIs
func (q Query) GetArtistsByURIs(uris ...string) ([]Artist, bool) {
	artists := make([]Artist, 0)
	ok := true

	for _, x := range uris {
		artist, thisOk := q.GetArtistByURI(x)
		ok = ok && thisOk
		artists = append(artists, artist)
	}
	return artists, ok
}

// Get multiple Artists by names
func (q Query) GetArtistsByNames(names ...string) ([]Artist, bool) {
	artists := make([]Artist, 0)
	ok := true

	for _, x := range names {
		artist, thisOk := q.GetArtistByName(x)
		ok = ok && thisOk
		artists = append(artists, artist)
	}
	return artists, ok
}

// Get multiple Playlists by URIs
func (q Query) GetPlaylistsByURIs(uris ...string) ([]Playlist, bool) {
	playlists := make([]Playlist, 0)
	ok := true

	for _, x := range uris {
		playlist, thisOk := q.GetPlaylistByURI(x)
		ok = ok && thisOk
		playlists = append(playlists, playlist)
	}
	return playlists, ok
}

// Get multiple Playlists by names
func (q Query) GetPlaylistsByNames(names ...string) ([]Playlist, bool) {
	playlists := make([]Playlist, 0)
	ok := true

	for _, x := range names {
		playlist, thisOk := q.GetPlaylistByName(x)
		ok = ok && thisOk
		playlists = append(playlists, playlist)
	}
	return playlists, ok
}
