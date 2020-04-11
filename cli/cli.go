package cli

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/gin-gonic/gin"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/sosedoff/x/wireguard-manager/assets"
	"github.com/sosedoff/x/wireguard-manager/controller"
	"github.com/sosedoff/x/wireguard-manager/model"
	"github.com/sosedoff/x/wireguard-manager/service"
	"github.com/sosedoff/x/wireguard-manager/store"
	"github.com/sosedoff/x/wireguard-manager/util"
)

func Run() {
	// TODO: how to set default on engine level
	setGinDefaults()

	var configPath string

	flag.StringVar(&configPath, "c", "", "Configuration file")
	flag.Parse()

	if configPath == "" {
		log.Fatal("config is required")
	}

	config, err := readConfig(configPath)
	if err != nil {
		log.Fatal("config error:", err)
	}

	datastore, err := store.Init(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := datastore.AutoMigrate(); err != nil {
		log.Fatal("automigrate error:", err)
	}

	ctl := controller.New(config.WireGuardBinPath, config.WireGuardPath, datastore)
	go ctl.Start()

	server, err := datastore.FindServer()
	if err != nil {
		log.Fatal(err)
	}
	if server == nil {
		log.Println("creating default server")
		s, err := configureDefaultServer(datastore)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("applying server config")
		if err := ctl.Apply(); err != nil {
			log.Fatal(err)
		}

		server = s
	}

	svc, err := service.New(&service.Config{
		AssetFS:         assets.Assets,
		Store:           datastore,
		Controller:      ctl,
		CookieName:      config.CookieName,
		CookieSecret:    config.CookieSecret,
		ClientID:        config.ClientID,
		ClientSecret:    config.ClientSecret,
		ClientDomain:    config.ClientDomain,
		ClientWhitelist: config.ClientWhitelist,
	})
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := fmt.Sprintf("%v:%v", "0.0.0.0", config.HTTPPort)

	log.Println("starting server on", listenAddr)
	if err := svc.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}

func setGinDefaults() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	log.SetFlags(log.LstdFlags)
}

func configureDefaultServer(store store.Store) (*model.Server, error) {
	endpoint, err := util.FetchPublicIP()
	if err != nil {
		return nil, err
	}

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	iface := "wg0"
	postUp := fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE; ip6tables -A FORWARD -i %s -j ACCEPT; ip6tables -t nat -A POSTROUTING -o eth0 -j MASQUERADE", iface, iface)
	postDown := fmt.Sprintf("iptables -D FORWARD -i %s -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE; ip6tables -D FORWARD -i %s -j ACCEPT; ip6tables -t nat -D POSTROUTING -o eth0 -j MASQUERADE", iface, iface)

	server := &model.Server{
		Name:       "private",
		Interface:  iface,
		Endpoint:   endpoint,
		PrivateKey: key.String(),
		PublicKey:  key.PublicKey().String(),
		IPV4Net:    "10.10.0.0/20",
		IPV4Addr:   "10.10.0.0/20",
		ListenPort: 51820,
		PostUp:     postUp,
		PostDown:   postDown,
	}

	err = store.CreateServer(server)
	return server, err
}
