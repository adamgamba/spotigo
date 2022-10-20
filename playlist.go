package spotigo

import (
	"fmt"
	"time"
)

// Playlist struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type Playlist struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalUrls  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		DisplayName  string `json:"display_name"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"owner"`
	PrimaryColor interface{} `json:"primary_color"`
	Public       bool        `json:"public"`
	SnapshotID   string      `json:"snapshot_id"`
	Tracks       struct {
		Href  string `json:"href"`
		Items []struct {
			AddedAt time.Time `json:"added_at"`
			AddedBy struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"added_by"`
			IsLocal        bool        `json:"is_local"`
			PrimaryColor   interface{} `json:"primary_color"`
			Track          Track       `json:"track"`
			VideoThumbnail struct {
				URL interface{} `json:"url"`
			} `json:"video_thumbnail"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// Get Playlist Name
func (p *Playlist) GetName() string {
	return p.Name
}

// Get Playlist URI
func (p *Playlist) GetURI() string {
	return p.ID
}

// Get number of Playlist followers
func (p *Playlist) GetNumFollowers() int {
	return p.Followers.Total
}

type tracksOfPlaylist struct {
	Href  string `json:"href"`
	Items []struct {
		AddedAt time.Time `json:"added_at"`
		AddedBy struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"added_by"`
		IsLocal        bool        `json:"is_local"`
		PrimaryColor   interface{} `json:"primary_color"`
		Track          Track       `json:"track"`
		VideoThumbnail struct {
			URL interface{} `json:"url"`
		} `json:"video_thumbnail"`
	} `json:"items"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

// Get all Track URIs for Playlist
func (p *Playlist) GetTrackURIs(u User) []string {
	uris := make([]string, 0)
	for _, track := range p.Tracks.Items {
		uris = append(uris, track.Track.ID)
	}

	// * added
	next := p.Tracks.Next
	fmt.Println("next val:", next)
	for next != "" {
		tracks := tracksOfPlaylist{}
		ok := u.sendGetRequest(next, &tracks)
		fmt.Println("gettrackURIs ok?", ok)

		next = tracks.Next
		fmt.Println("next val:", next)

	}

	return uris
}

// Get all Artist URIs for Playlist
func (p *Playlist) GetArtistURIs() []string {
	uris := make([]string, 0)
	for _, x := range p.Tracks.Items {
		for _, artist := range x.Track.Artists {
			uris = append(uris, artist.ID)
		}
	}
	return uris
}

// Get all Tracks on a Playlist
func (p *Playlist) GetTracks() []Track {
	uris := make([]Track, 0)
	for _, x := range p.Tracks.Items {
		uris = append(uris, x.Track)
	}
	return uris
}
