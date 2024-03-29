package spotigo

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Get a User's Saved Tracks
// limit sets the number of tracks to return
// If getAll is true, limit is disregarded
func (u *User) GetSavedTracks(getAll bool, limit int) ([]Track, bool) {
	const MAX_LIMIT = 50
	tracks := make([]Track, 0)
	ok := true
	offset := 0

	if limit < 0 {
		return nil, false
	}
	if getAll {
		numTracks, ok2 := u.GetNumSavedTracks()
		ok = ok && ok2
		limit = numTracks
	}

	// Spotify will only return 50 tracks at a time
	// If limit > 50, loop through 50 tracks per request
	for i := 0; limit > MAX_LIMIT; i++ {

		reqURL := u.baseURL + "me/tracks?limit=" + fmt.Sprint(MAX_LIMIT) + "&offset=" + fmt.Sprint(offset)
		savedTracks := savedTracks{}

		ok = u.sendGetRequest(reqURL, &savedTracks) && ok

		for _, x := range savedTracks.Items {
			tracks = append(tracks, x.Track)
		}

		offset += MAX_LIMIT
		limit -= MAX_LIMIT

	}

	// Format request URL
	reqURL := u.baseURL + "me/tracks?limit=" + fmt.Sprint(limit) + "&offset=" + fmt.Sprint(offset)
	savedTracks := savedTracks{}

	// Send request, store results in savedTracks
	ok = u.sendGetRequest(reqURL, &savedTracks) && ok

	for _, x := range savedTracks.Items {
		tracks = append(tracks, x.Track)
	}

	return tracks, ok
}

// Get the number of tracks a User has saved
func (u *User) GetNumSavedTracks() (int, bool) {
	reqURL := u.baseURL + "me/tracks?limit=1"
	savedTracks := savedTracks{}
	ok := u.sendGetRequest(reqURL, &savedTracks)
	return savedTracks.Total, ok
}

// Saved Tracks struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type savedTracks struct {
	Href  string `json:"href"`
	Items []struct {
		AddedAt time.Time `json:"added_at"`
		Track   Track     `json:"track"`
	} `json:"items"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

// Get a User's Saved Albums
// limit sets the number of albums to return
// If getAll is true, limit is disregarded
func (u *User) GetSavedAlbums(getAll bool, limit int) ([]Album, bool) {
	const MAX_LIMIT = 50
	albums := make([]Album, 0)
	ok := true
	offset := 0

	if limit < 0 {
		return nil, false
	}
	if getAll {
		numAlbums, ok2 := u.GetNumSavedAlbums()
		ok = ok && ok2
		limit = numAlbums
	}

	// Spotify will only return 50 albums at a time
	// If limit > 50, loop through 50 albums per request
	for i := 0; limit > MAX_LIMIT; i++ {

		reqURL := u.baseURL + "me/albums?limit=" + fmt.Sprint(MAX_LIMIT) + "&offset=" + fmt.Sprint(offset)
		savedAlbums := savedAlbums{}

		ok = u.sendGetRequest(reqURL, &savedAlbums) && ok

		for _, x := range savedAlbums.Items {
			albums = append(albums, x.Album)
		}

		offset += MAX_LIMIT
		limit -= MAX_LIMIT

	}

	reqURL := u.baseURL + "me/albums?limit=" + fmt.Sprint(limit) + "&offset=" + fmt.Sprint(offset)
	savedAlbums := savedAlbums{}

	ok = u.sendGetRequest(reqURL, &savedAlbums) && ok

	for _, x := range savedAlbums.Items {
		albums = append(albums, x.Album)
	}

	return albums, ok
}
func (u *User) GetNumSavedAlbums() (int, bool) {
	reqURL := u.baseURL + "me/albums?limit=1"
	savedAlbums := savedAlbums{}
	ok := u.sendGetRequest(reqURL, &savedAlbums)
	return savedAlbums.Total, ok
}

// Saved Albums struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type savedAlbums struct {
	Href  string `json:"href"`
	Items []struct {
		AddedAt time.Time `json:"added_at"`
		Album   Album     `json:"album"`
	} `json:"items"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

// Get a User's Saved Playlists
// limit sets the number of playlists to return
// If getAll is true, limit is disregarded
func (u *User) GetSavedPlaylists(getAll bool, limit int) ([]Playlist, bool) {
	const MAX_LIMIT = 50
	playlists := make([]Playlist, 0)
	ok := true
	offset := 0

	if limit < 0 {
		return nil, false
	}
	if getAll {
		numPlaylists, ok2 := u.GetNumSavedPlaylists()
		ok = ok && ok2
		limit = numPlaylists
	}

	// Spotify will only return 50 playlists at a time
	// If limit > 50, loop through 50 playlists per request
	for i := 0; limit > MAX_LIMIT; i++ {

		reqURL := u.baseURL + "me/playlists?limit=" + fmt.Sprint(MAX_LIMIT) + "&offset=" + fmt.Sprint(offset)
		savedPlaylists := savedPlaylists{}

		ok = u.sendGetRequest(reqURL, &savedPlaylists) && ok

		for _, x := range savedPlaylists.Items {
			playlists = append(playlists, x)
		}

		offset += MAX_LIMIT
		limit -= MAX_LIMIT

	}

	reqURL := u.baseURL + "me/playlists?limit=" + fmt.Sprint(limit) + "&offset=" + fmt.Sprint(offset)
	savedPlaylists := savedPlaylists{}

	ok = u.sendGetRequest(reqURL, &savedPlaylists) && ok

	for _, x := range savedPlaylists.Items {
		playlists = append(playlists, x)
	}

	return playlists, ok
}
func (u *User) GetNumSavedPlaylists() (int, bool) {
	reqURL := u.baseURL + "me/playlists?limit=1"
	savedPlaylists := savedPlaylists{}
	ok := u.sendGetRequest(reqURL, &savedPlaylists)
	return savedPlaylists.Total, ok
}

// Saved Playlists struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type savedPlaylists struct {
	Href     string      `json:"href"`
	Items    []Playlist  `json:"items"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

// Get a User's Saved Artists
// limit sets the number of artists to return
// If getAll is true, limit is disregarded
func (u *User) GetFollowedArtists(getAll bool, limit int) ([]Artist, bool) {
	const MAX_LIMIT = 50
	artists := make([]Artist, 0)
	ok := true
	after := ""

	if limit < 0 {
		return nil, false
	}
	if getAll {
		numArtists, ok2 := u.GetNumFollowedArtists()
		ok = ok && ok2
		limit = numArtists
	}

	// Spotify will only return 50 artists at a time
	// If limit > 50, loop through 50 artists per request
	for i := 0; limit > MAX_LIMIT; i++ {

		reqURL := u.baseURL + "me/following?type=artist&limit=" + fmt.Sprint(MAX_LIMIT)
		if after != "" {
			reqURL += "&after=" + after
		}

		followedArtists := followedArtists{}

		ok = u.sendGetRequest(reqURL, &followedArtists) && ok

		for _, x := range followedArtists.Artists.Items {
			artists = append(artists, x)
		}
		after = followedArtists.Artists.Cursors.After

		if i == 0 && getAll {
			limit = followedArtists.Artists.Total
		}

		// offset += MAX_LIMIT
		limit -= MAX_LIMIT
	}

	reqURL := u.baseURL + "me/following?type=artist&limit=" + fmt.Sprint(limit)
	if after != "" {
		reqURL += "&after=" + after
	}

	followedArtists := followedArtists{}

	ok = u.sendGetRequest(reqURL, &followedArtists) && ok

	for _, x := range followedArtists.Artists.Items {
		artists = append(artists, x)
	}
	after = followedArtists.Artists.Cursors.After

	return artists, ok
}
func (u *User) GetNumFollowedArtists() (int, bool) {
	reqURL := u.baseURL + "me/following?type=artist&limit=1"
	followedArtists := followedArtists{}
	ok := u.sendGetRequest(reqURL, &followedArtists)
	return followedArtists.Artists.Total, ok
}

// Followed Artists struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type followedArtists struct {
	Artists struct {
		Items   []Artist `json:"items"`
		Next    string   `json:"next"`
		Total   int      `json:"total"`
		Cursors struct {
			After string `json:"after"`
		} `json:"cursors"`
		Limit int    `json:"limit"`
		Href  string `json:"href"`
	} `json:"artists"`
}

// Get audio features for a track
// Examples of audio features include Danceability, Energy, Valence
func (u *User) GetTrackAudioFeatures(q Query, i interface{}) (AudioFeatures, bool) {
	uri := ""
	ok := true

	switch v := i.(type) {
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
			uri = track.ID
		}
	case Track:
		uri = v.ID
	default:
		fmt.Println("Invalid Type")
		ok = false
	}
	reqURL := u.baseURL + "audio-features/" + uri

	audioFeatures := AudioFeatures{}
	ok = u.sendGetRequest(reqURL, &audioFeatures) && ok
	return audioFeatures, ok
}

// Get the audio analysis of a track
// Returned values include duration, number of samples, and loudness
func (u *User) GetTrackAudioAnalysis(q Query, i interface{}) (AudioAnalysis, bool) {
	uri := ""
	ok := true

	switch v := i.(type) {
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
			uri = track.ID
		}
	case Track:
		uri = v.ID
	default:
		ok = false
		fmt.Println("Invalid Type")
	}
	reqURL := u.baseURL + "audio-analysis/" + uri
	fmt.Println("requrl", reqURL)

	audioAnalysis := AudioAnalysis{}
	ok = u.sendGetRequest(reqURL, &audioAnalysis) && ok
	return audioAnalysis, ok
}

// Audio Features struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type AudioFeatures struct {
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Type             string  `json:"type"`
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	TrackHref        string  `json:"track_href"`
	AnalysisURL      string  `json:"analysis_url"`
	DurationMs       int     `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

// Audio Analysis struct- maps to Spotify JSON response format by tag `json: "var_name"`
// Struct generated by putting Spotify JSON data into JSON to Go struct generator at:
// https://mholt.github.io/json-to-go/
type AudioAnalysis struct {
	Meta struct {
		AnalyzerVersion string  `json:"analyzer_version"`
		Platform        string  `json:"platform"`
		DetailedStatus  string  `json:"detailed_status"`
		StatusCode      int     `json:"status_code"`
		Timestamp       int     `json:"timestamp"`
		AnalysisTime    float64 `json:"analysis_time"`
		InputProcess    string  `json:"input_process"`
	} `json:"meta"`
	Track struct {
		NumSamples              int     `json:"num_samples"`
		Duration                float64 `json:"duration"`
		SampleMd5               string  `json:"sample_md5"`
		OffsetSeconds           int     `json:"offset_seconds"`
		WindowSeconds           int     `json:"window_seconds"`
		AnalysisSampleRate      int     `json:"analysis_sample_rate"`
		AnalysisChannels        int     `json:"analysis_channels"`
		EndOfFadeIn             float64 `json:"end_of_fade_in"`
		StartOfFadeOut          float64 `json:"start_of_fade_out"`
		Loudness                float64 `json:"loudness"`
		Tempo                   float64 `json:"tempo"`
		TempoConfidence         float64 `json:"tempo_confidence"`
		TimeSignature           int     `json:"time_signature"`
		TimeSignatureConfidence float64 `json:"time_signature_confidence"`
		Key                     int     `json:"key"`
		KeyConfidence           float64 `json:"key_confidence"`
		Mode                    int     `json:"mode"`
		ModeConfidence          float64 `json:"mode_confidence"`
		Codestring              string  `json:"codestring"`
		CodeVersion             float64 `json:"code_version"`
		Echoprintstring         string  `json:"echoprintstring"`
		EchoprintVersion        float64 `json:"echoprint_version"`
		Synchstring             string  `json:"synchstring"`
		SynchVersion            float64 `json:"synch_version"`
		Rhythmstring            string  `json:"rhythmstring"`
		RhythmVersion           float64 `json:"rhythm_version"`
	} `json:"track"`
	Bars []struct {
		Start      float64 `json:"start"`
		Duration   float64 `json:"duration"`
		Confidence float64 `json:"confidence"`
	} `json:"bars"`
	Beats []struct {
		Start      float64 `json:"start"`
		Duration   float64 `json:"duration"`
		Confidence float64 `json:"confidence"`
	} `json:"beats"`
	Sections []struct {
		Start                   float64 `json:"start"`
		Duration                float64 `json:"duration"`
		Confidence              float64 `json:"confidence"`
		Loudness                float64 `json:"loudness"`
		Tempo                   float64 `json:"tempo"`
		TempoConfidence         float64 `json:"tempo_confidence"`
		Key                     int     `json:"key"`
		KeyConfidence           float64 `json:"key_confidence"`
		Mode                    int     `json:"mode"`
		ModeConfidence          float64 `json:"mode_confidence"`
		TimeSignature           int     `json:"time_signature"`
		TimeSignatureConfidence float64 `json:"time_signature_confidence"`
	} `json:"sections"`
	Segments []struct {
		Start           float64   `json:"start"`
		Duration        float64   `json:"duration"`
		Confidence      float64   `json:"confidence"`
		LoudnessStart   float64   `json:"loudness_start"`
		LoudnessMaxTime float64   `json:"loudness_max_time"`
		LoudnessMax     float64   `json:"loudness_max"`
		LoudnessEnd     float64   `json:"loudness_end"`
		Pitches         []float64 `json:"pitches"`
		Timbre          []float64 `json:"timbre"`
	} `json:"segments"`
	Tatums []struct {
		Start      float64 `json:"start"`
		Duration   float64 `json:"duration"`
		Confidence float64 `json:"confidence"`
	} `json:"tatums"`
}

func (u *User) GetCurrentProfile() (Profile, bool) {
	reqURL := u.baseURL + "me"

	profile := Profile{}
	ok := u.sendGetRequest(reqURL, &profile)
	return profile, ok
}
