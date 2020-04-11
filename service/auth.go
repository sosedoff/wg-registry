package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

var (
	authStartPath = "/auth/google"
)

type GoogleAuth struct {
	UID    string
	Email  string
	Domain string
	Name   string
}

func domainFromEmail(val string) string {
	chunks := strings.Split(val, "@")
	return chunks[len(chunks)-1]
}

type AuthConfig struct {
	oauth     *oauth2.Config
	domain    string
	whitelist []string
}

func googleAuthFromResponse(resp *http.Response) (*GoogleAuth, error) {
	auth := &GoogleAuth{}

	fields := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&fields); err != nil {
		return nil, err
	}

	auth.UID = fields["sub"].(string)
	auth.Email = fields["email"].(string)
	auth.Name = fields["name"].(string)
	auth.Domain = domainFromEmail(auth.Email)

	return auth, nil
}
