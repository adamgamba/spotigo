package spotigo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Device contains information about a device that a user can play music on
// Source: https://github.com/zmb3/spotify/
type Device struct {
	// ID of the device. This may be empty.
	ID string `json:"id"`
	// Active If this device is the currently active device.
	Active bool `json:"is_active"`
	// Restricted Whether controlling this device is restricted. At present if
	// this is "true" then no Web API commands will be accepted by this device.
	Restricted bool `json:"is_restricted"`
	// Name The name of the device.
	Name string `json:"name"`
	// Type of device, such as "Computer", "Smartphone" or "Speaker".
	Type string `json:"type"`
	// Volume The current volume in percent.
	Volume int `json:"volume_percent"`
}

// Return a user's available playback devices
func (u *User) GetPlaybackDevices() ([]Device, bool) {
	var result struct {
		Devices []Device `json:"devices"`
	}

	reqURL := u.baseURL + "me/player/devices"

	ok := u.sendGetRequest(reqURL, &result)

	return result.Devices, ok
}

// Transfer playback between devices
func (u *User) TransferPlayback(device interface{}, play bool) bool {

	deviceID := ""
	switch v := device.(type) {
	case string:
		deviceID = v
	case Device:
		deviceID = string(v.ID)
	// Base Case: Invalid Type
	default:
		return false
	}

	// Source: https://github.com/zmb3/spotify/
	reqData := struct {
		DeviceID []string `json:"device_ids"`
		Play     bool     `json:"play"`
	}{
		DeviceID: []string{deviceID},
		Play:     play,
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reqData)
	/// End Source

	if err != nil {
		return false
	}

	req, err := http.NewRequest(http.MethodPut, u.baseURL+"me/player", buf)
	if err != nil {
		return false
	}
	err = u.execute(req, nil, http.StatusNoContent)
	if err != nil {
		return false
	}

	return true
}
