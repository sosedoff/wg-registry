package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sosedoff/wg-registry/util"

	"github.com/sosedoff/wg-registry/service"
)

type FileConfig struct {
	Store           string                     `json:"store"`
	CookieName      string                     `json:"cookie_name"`
	CookieSecret    string                     `json:"cookie_secret"`
	ClientID        string                     `json:"client_id"`
	ClientSecret    string                     `json:"client_secret"`
	ClientDomain    string                     `json:"client_domain"`
	ClientWhitelist []string                   `json:"client_whitelist"`
	WGDir           string                     `json:"wg_dir"`
	WGPath          string                     `json:"wg_path"`
	WGQuickPath     string                     `json:"wg_quick_path"`
	HTTPPort        int                        `json:"http_port"`
	HTTPSPort       int                        `json:"https_port"`
	Debug           bool                       `json:"debug"`
	LetsEncrypt     *service.LetsEncryptConfig `json:"letsencrypt"`
}

func readConfig(path string) (*FileConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &FileConfig{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.ClientID == "" {
		return nil, errors.New("client_id is not set")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("client_secret is not set")
	}
	if config.ClientDomain == "" {
		return nil, errors.New("client_domain is not set")
	}
	if config.Store == "" {
		config.Store = "/etc/wg-registry/data.json"
	}
	if config.CookieName == "" {
		config.CookieName = "wg-registry"
	}
	if config.CookieSecret == "" {
		hex, err := util.RandomHex(32)
		if err == nil {
			config.CookieSecret = hex
		} else {
			config.CookieSecret = fmt.Sprintf("%v", time.Now().UnixNano())
		}
	}
	if config.HTTPPort == 0 {
		config.HTTPPort = 80
	}
	if config.HTTPSPort == 0 {
		config.HTTPSPort = 443
	}
	if config.WGDir == "" {
		config.WGDir = "/etc/wireguard"
	}
	if config.WGPath == "" {
		config.WGPath = "wg"
	}
	if config.WGQuickPath == "" {
		config.WGQuickPath = "wg-quick"
	}

	return config, nil
}
