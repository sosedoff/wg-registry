package util

import (
	"encoding/json"
	"net/http"
)

// FetchPublicIP returns a public IP address of the machine
func FetchPublicIP() (string, error) {
	resp, err := http.Get("https://ipinfo.io/what-is-my-ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reply := &struct {
		IP string `json:"ip"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(reply); err != nil {
		return "", err
	}

	return reply.IP, nil
}
