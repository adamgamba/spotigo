package spotigo

import (
	"encoding/json"
	"log"
	"net/http"
)

// Save tracks for User
func (u *User) SaveTracks(q Query, i ...interface{}) bool {
	return u.modifySavedTracks(q, true, i...)
}

// Unsave tracks for User
func (u *User) UnsaveTracks(q Query, i ...interface{}) bool {
	return u.modifySavedTracks(q, false, i...)
}

// Execute saving/unsaving of tracks
func (u *User) modifySavedTracks(q Query, save bool, i ...interface{}) bool {
	uris := make([]string, 0)
	ok := true

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "track")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Tracks.Total > 0 {
				track := searchResult.Tracks.Items[0]
				uri := track.ID
				uris = append(uris, uri)
			}
		case Track:
			uri := v.ID
			uris = append(uris, uri)
		// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return ok
		}
	}

	reqURL := u.baseURL + "me/tracks?ids="
	for _, uri := range uris {
		reqURL += uri + ","
	}
	reqURL = reqURL[:len(reqURL)-1]

	method := ""
	if save {
		method = http.MethodPut
	} else {
		method = http.MethodDelete
	}

	ok = u.sendRequest(method, reqURL) && ok

	return ok
}

// Follow Artists for User
func (u *User) FollowArtists(q Query, i ...interface{}) bool {
	return u.modifyFollowees(q, true, true, i...)
}

// Unfollow Artists for User
func (u *User) UnfollowArtists(q Query, i ...interface{}) bool {
	return u.modifyFollowees(q, false, true, i...)
}

// Follow Users for User
func (u *User) FollowUsers(q Query, i ...interface{}) bool {
	return u.modifyFollowees(q, true, false, i...)
}

// Unfollow Users for User
func (u *User) UnfollowUsers(q Query, i ...interface{}) bool {
	return u.modifyFollowees(q, false, false, i...)
}

// Execute following/unfollowing of artists or users
func (u *User) modifyFollowees(q Query, follow bool, artist bool, i ...interface{}) bool {
	uris := make([]string, 0)
	ok := true

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "artist")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Artists.Total > 0 {
				artist := searchResult.Artists.Items[0]
				uri := artist.ID
				uris = append(uris, uri)
			}
		case Artist:
			uri := v.ID
			uris = append(uris, uri)
		// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return ok
		}
	}

	reqURL := u.baseURL + "me/following?ids="
	for _, uri := range uris {
		reqURL += uri + ","
	}
	reqURL = reqURL[:len(reqURL)-1]

	qType := ""
	if artist {
		qType = "artist"
	} else {
		qType = "user"
	}
	reqURL += "&type=" + qType

	method := ""
	if follow {
		method = http.MethodPut
	} else {
		method = http.MethodDelete
	}

	ok = u.sendRequest(method, reqURL) && ok

	return ok
}

// Save Playlists for User
func (u *User) SavePlaylists(q Query, i ...interface{}) bool {
	return u.modifySavedPlaylists(q, true, i...)
}

// Unsave Playlists for User
func (u *User) UnsavePlaylists(q Query, i ...interface{}) bool {
	return u.modifySavedPlaylists(q, false, i...)
}

// Execute saving/unsaving of playlists
func (u *User) modifySavedPlaylists(q Query, save bool, i ...interface{}) bool {

	method := ""
	if save {
		method = http.MethodPut
	} else {
		method = http.MethodDelete
	}

	uri := ""
	ok := true

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "playlist")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Playlists.Total > 0 {
				playlist := searchResult.Playlists.Items[0]
				uri = playlist.ID
			}
		case Playlist:
			uri = v.ID
		// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return ok
		}

		reqURL := u.baseURL + "playlists/" + uri + "/followers"
		ok = u.sendRequest(method, reqURL) && ok

	}

	return ok
}

// Save Albums for User
func (u *User) SaveAlbums(q Query, i ...interface{}) bool {
	return u.modifySavedAlbums(q, true, i...)
}

// Unsave Albums for User
func (u *User) UnsaveAlbums(q Query, i ...interface{}) bool {
	return u.modifySavedAlbums(q, false, i...)
}

// Execute saving/unsaving of albums
func (u *User) modifySavedAlbums(q Query, save bool, i ...interface{}) bool {
	uris := make([]string, 0)
	ok := true

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "album")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Albums.Total > 0 {
				album := searchResult.Albums.Items[0]
				uri := album.ID
				uris = append(uris, uri)
			}
		case Album:
			uri := v.ID
			uris = append(uris, uri)
			// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return ok
		}
	}

	reqURL := u.baseURL + "me/albums?ids="
	for _, uri := range uris {
		reqURL += uri + ","
	}
	reqURL = reqURL[:len(reqURL)-1]

	method := ""
	if save {
		method = http.MethodPut
	} else {
		method = http.MethodDelete
	}

	ok = u.sendRequest(method, reqURL) && ok

	return ok
}

// Search Result struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type searchResult struct {
	Tracks struct {
		Href     string  `json:"href"`
		Items    []Track `json:"items"`
		Limit    int     `json:"limit"`
		Next     string  `json:"next"`
		Offset   int     `json:"offset"`
		Previous string  `json:"previous"`
		Total    int     `json:"total"`
	} `json:"tracks"`
	Artists struct {
		Href     string   `json:"href"`
		Items    []Artist `json:"items"`
		Limit    int      `json:"limit"`
		Next     string   `json:"next"`
		Offset   int      `json:"offset"`
		Previous string   `json:"previous"`
		Total    int      `json:"total"`
	} `json:"artists"`
	Albums struct {
		Href     string  `json:"href"`
		Items    []Album `json:"items"`
		Limit    int     `json:"limit"`
		Next     string  `json:"next"`
		Offset   int     `json:"offset"`
		Previous string  `json:"previous"`
		Total    int     `json:"total"`
	} `json:"albums"`
	Playlists struct {
		Href     string     `json:"href"`
		Items    []Playlist `json:"items"`
		Limit    int        `json:"limit"`
		Next     string     `json:"next"`
		Offset   int        `json:"offset"`
		Previous string     `json:"previous"`
		Total    int        `json:"total"`
	} `json:"playlists"`
}

// Check if a User is following a set of artists
// Returns a list of booleans corresponding to whether that artist in the
// parameter list is followed by the User
func (u *User) DoesFollowArtists(q Query, i ...interface{}) ([]bool, bool) {
	ok := true
	uris := make([]string, 0)

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "artist")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Artists.Total > 0 {
				playlist := searchResult.Artists.Items[0]
				uris = append(uris, playlist.ID)
			}
		case Artist:
			uris = append(uris, v.ID)
		// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return make([]bool, 0), ok
		}
	}

	reqURL := u.baseURL + "me/following/contains?ids="
	for _, uri := range uris {
		reqURL += uri + ","
	}
	reqURL = reqURL[:len(reqURL)-1]
	reqURL += "&type=artist"

	bools := make([]bool, 0)
	ok = u.sendGetRequest(reqURL, &bools) && ok

	return bools, ok
}

// Check if a User has saved a set of tracks
// Returns a list of booleans corresponding to whether that track in the
// parameter list is followed by the User
func (u *User) HasSavedTracks(q Query, i ...interface{}) ([]bool, bool) {
	ok := true
	uris := make([]string, 0)

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "track")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Tracks.Total > 0 {
				playlist := searchResult.Tracks.Items[0]
				uris = append(uris, playlist.ID)
			}
		case Track:
			uris = append(uris, v.ID)
		// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return make([]bool, 0), ok
		}
	}

	reqURL := u.baseURL + "me/tracks/contains?ids="
	for _, uri := range uris {
		reqURL += uri + ","
	}
	reqURL = reqURL[:len(reqURL)-1]

	bools := make([]bool, 0)
	ok = u.sendGetRequest(reqURL, &bools) && ok

	return bools, ok
}

// Check if a User has saved a set of albums
// Returns a list of booleans corresponding to whether that album in the
// parameter list is followed by the User
func (u *User) HasSavedAlbums(q Query, i ...interface{}) ([]bool, bool) {
	ok := true
	uris := make([]string, 0)

	for _, item := range i {
		switch v := item.(type) {
		case string:
			bytes, searchErr := q.search(v, "album")
			if searchErr != nil {
				ok = false
			}

			searchResult := searchResult{}
			jsonErr := json.Unmarshal(bytes, &searchResult)

			if jsonErr != nil {
				log.Fatal("JSON Error:", jsonErr)
			}

			if searchResult.Albums.Total > 0 {
				playlist := searchResult.Albums.Items[0]
				uris = append(uris, playlist.ID)
			}
		case Album:
			uris = append(uris, v.ID)
		// Invalid Type
		default:
			ok = false
		}

		if !ok {
			return make([]bool, 0), ok
		}
	}

	reqURL := u.baseURL + "me/albums/contains?ids="
	for _, uri := range uris {
		reqURL += uri + ","
	}
	reqURL = reqURL[:len(reqURL)-1]

	bools := make([]bool, 0)
	ok = u.sendGetRequest(reqURL, &bools) && ok

	return bools, ok
}
