package service

import (
	"github.com/jessevdk/go-assets"
	"github.com/sosedoff/x/wireguard-manager/controller"
	"github.com/sosedoff/x/wireguard-manager/store"
)

type Config struct {
	AssetFS         *assets.FileSystem
	Store           store.Store
	Controller      *controller.Controller
	CookieName      string
	CookieSecret    string
	ClientID        string
	ClientSecret    string
	ClientDomain    string
	ClientWhitelist []string
}
