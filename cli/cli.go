package cli

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/sosedoff/wg-registry/assets"
	"github.com/sosedoff/wg-registry/controller"
	"github.com/sosedoff/wg-registry/service"
	"github.com/sosedoff/wg-registry/store"
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
	if server != nil {
		log.Println("applying server config")
		if err := ctl.Apply(); err != nil {
			log.Fatal(err)
		}
	}

	httpsEnabled := config.LetsEncrypt != nil && config.LetsEncrypt.Enabled == true

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
		ForceHTTPS:      httpsEnabled,
	})
	if err != nil {
		log.Fatal(err)
	}

	if httpsEnabled {
		certManager, err := service.NewCertManager(config.LetsEncrypt)
		if err != nil {
			log.Fatal(err)
		}

		server := &http.Server{
			Addr:      fmt.Sprintf("%v:%v", "0.0.0.0", config.HTTPSPort),
			Handler:   svc,
			TLSConfig: certManager.TLSConfig(),
		}

		go func() {
			log.Println("starting https listener on", server.Addr)
			if err := server.ListenAndServeTLS("", ""); err != nil {
				log.Fatal("https listener error:", err)
			}
		}()
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
