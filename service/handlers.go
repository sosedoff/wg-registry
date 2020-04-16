package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"golang.org/x/oauth2"
	"xojoc.pw/useragent"

	"github.com/sosedoff/wg-registry/generate"
	"github.com/sosedoff/wg-registry/model"
	"github.com/sosedoff/wg-registry/util"
)

func New(config *Config) (*gin.Engine, error) {
	router := gin.New()

	// Defaults
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// HTML templates
	tpl, err := loadTemplate(config.AssetFS)
	if err != nil {
		return nil, err
	}
	router.SetHTMLTemplate(tpl)

	// Middlware
	if config.ForceHTTPS {
		router.Use(forceHTTPS())
	}

	router.Use(setAuth(config))
	router.Use(setStore(config))
	router.Use(setController(config))
	router.Use(rejectBots())
	router.Use(cookieSession(config.CookieName, config.CookieSecret))

	// User routes
	router.GET("/", requireUser, handleIndex)
	router.POST("/devices", requireUser, handleCreateDevice)
	router.GET("/devices/:id", requireUser, requireDevice, handleDevice)
	router.GET("/devices/:id/config", requireUser, requireDevice, handleDeviceConfig)
	router.GET("/devices/:id/delete", requireUser, requireDevice, handleDeleteDevice)

	// Authentication
	router.GET("/auth/google", handleAuthIndex)
	router.GET("/auth/google/start", handleAuthStart)
	router.GET("/auth/google/callback", handleAuthCallback)
	router.GET("/auth/signout", requireUser, handleSignout)

	// Admin routes
	admin := router.Group("/admin", requireUser, requireAdmin)
	admin.GET("", handleAdminIndex)
	admin.GET("/server", handleAdminServer)
	admin.POST("/server", handleAdminServer)
	admin.GET("/server/restart", handleAdminRestartServer)

	// Everything else
	router.GET("/static/:file", serveStaticAsset(config.AssetFS))

	return router, nil
}

func handleAuthIndex(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"startPath": "/auth/google/start",
	})
}

func handleAuthStart(c *gin.Context) {
	redirectURL := makeRedirectURL(c.Request, "/auth/google/callback")
	conf := getAuth(c).oauth
	conf.RedirectURL = redirectURL

	url := conf.AuthCodeURL(redirectURL)
	c.Redirect(302, url)
}

func handleAuthCallback(c *gin.Context) {
	auth := getAuth(c)
	conf := auth.oauth
	conf.RedirectURL = c.Query("state")

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		htmlError(c, err)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		htmlError(c, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		htmlError(c, "Google user info request failed")
		return
	}

	info, err := googleAuthFromResponse(resp)
	if err != nil {
		htmlError(c, err)
		return
	}

	if info.Domain != auth.domain {
		whitelisted := false
		for _, email := range auth.whitelist {
			if email == info.Email {
				whitelisted = true
				break
			}
		}
		if !whitelisted {
			htmlError(c, "Email address is not permitted")
			return
		}
	}

	store := getStore(c)
	user, err := store.FindUserByEmail(info.Email)
	if err != nil {
		htmlError(c, err)
		return
	}
	if user == nil {
		user = &model.User{
			Email:   info.Email,
			Name:    info.Name,
			Role:    roleUser,
			Enabled: true,
		}

		count, err := store.UserCount()
		if err != nil {
			htmlError(c, err)
			return
		}
		if count == 0 {
			user.Role = roleAdmin
		}

		if err := store.CreateUser(user); err != nil {
			htmlError(c, err)
			return
		}

	}
	session := getSession(c)
	session.Set("uid", user.ID)
	session.Save()

	c.Redirect(302, "/")
}

func handleSignout(c *gin.Context) {
	session := getSession(c)
	session.Delete("uid")
	session.Save()

	c.Redirect(302, "/")
}

func handleIndex(c *gin.Context) {
	store := getStore(c)
	user := getUser(c)

	server, err := store.FindServer()
	if err != nil {
		htmlError(c, err)
		return
	}
	if server == nil && user.Role == roleAdmin {
		c.Redirect(302, "/admin/server")
		c.Abort()
		return
	}

	devices, err := store.GetDevicesByUser(user.ID)
	if err != nil {
		htmlError(c, err)
		return
	}

	name := ""
	if agent := useragent.Parse(c.Request.UserAgent()); agent != nil {
		name = fmt.Sprintf("%s", agent.OS)
	}

	c.HTML(200, "index.html", gin.H{
		"user":          user,
		"devices":       devices,
		"newDeviceName": name,
	})
}

func handleDevice(c *gin.Context) {
	c.HTML(200, "device.html", gin.H{
		"device": getDevice(c),
	})
}

func handleCreateDevice(c *gin.Context) {
	store := getStore(c)
	user := getUser(c)
	ctl := getController(c)

	server, err := store.FindServer()
	if err != nil {
		htmlError(c, err)
		return
	}
	if server == nil {
		htmlError(c, "Server does not exist")
		return
	}

	device := &model.Device{
		UserID:              user.ID,
		Name:                c.Request.FormValue("name"),
		OS:                  c.Request.FormValue("os"),
		Enabled:             true,
		PersistentKeepalive: 60,
	}

	if err := store.CreateDevice(server, device); err != nil {
		htmlError(c, err)
		return
	}

	// Reload wireguard configuration
	if ctl != nil {
		if err := ctl.Apply(false); err != nil {
			log.Println("wireguard reload error:", err)
		}
	}

	c.Redirect(302, fmt.Sprintf("/devices/%d", device.ID))
}

func handleDeviceConfig(c *gin.Context) {
	store := getStore(c)
	device := getDevice(c)

	server, err := store.FindServer()
	if err != nil {
		badRequest(c, err)
		return
	}

	config, err := generate.ClientConfig(store, device, server)
	if err != nil {
		badRequest(c, err)
		return
	}

	if c.Query("qr") == "1" {
		qrdata, err := qrcode.Encode(string(config), qrcode.Medium, 400)
		if err != nil {
			badRequest(c, err)
			return
		}
		c.Data(200, "image/png", qrdata)
		return
	}

	if c.Query("view") == "" {
		filename := fmt.Sprintf("%s.conf", server.Name)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	}

	c.String(200, "%s", config)
}

func handleDeleteDevice(c *gin.Context) {
	user := getUser(c)
	store := getStore(c)
	ctl := getController(c)

	device, err := store.FindUserDevice(user, c.Param("id"))
	if err != nil {
		htmlError(c, err)
		return
	}
	if device == nil {
		c.Redirect(302, "/")
		return
	}

	if err := store.DeleteUserDevice(user, device); err != nil {
		htmlError(c, err)
		return
	}

	if ctl != nil {
		if err := ctl.Apply(false); err != nil {
			log.Println("wireguard reload error:", err)
		}
	}

	c.Redirect(302, "/")
}

func handleAdminIndex(c *gin.Context) {
	store := getStore(c)

	users, err := store.AllUsers()
	if err != nil {
		htmlError(c, err)
		return
	}

	devices, err := store.AllDevices()
	if err != nil {
		htmlError(c, err)
		return
	}

	server, err := store.FindServer()
	if err != nil {
		htmlError(c, err)
		return
	}

	config, err := generate.ServerConfig(store, server)
	if err != nil {
		htmlError(c, err)
		return
	}

	peers, err := fetchWireGuardPeers(server.Interface)
	if err != nil {
		log.Println("cant fetch peers:", err)
	}
	if peers != nil {
		for _, p := range peers {
			key := p.PublicKey.String()
			for _, d := range devices {
				if d.PublicKey == key {
					d.SetPeerInfo(&p)
					break
				}
			}
		}
	}

	c.HTML(200, "admin_index.html", gin.H{
		"devices": devices,
		"users":   users,
		"server":  server,
		"config":  string(config),
	})
}

func handleAdminServer(c *gin.Context) {
	store := getStore(c)
	ctl := getController(c)

	var err error
	server, err := store.FindServer()
	if err != nil {
		htmlError(c, err)
		return
	}
	if server == nil {
		server = model.ServerWithDefaults()
	}

	if c.Request.Method == http.MethodGet && server.ID == 0 {
		ip, err := util.FetchPublicIP()
		if err != nil {
			htmlError(c, err)
			return
		}
		server.Endpoint = ip
	} else if c.Request.Method == http.MethodPost {
		err = util.ErrChain(
			func() error { return c.ShouldBind(&server) },
			func() error { return server.Validate() },
			func() error { return server.AssignPrivateKey() },
			func() error { return store.SaveServer(server) },
			func() error {
				if ctl != nil {
					return ctl.Apply(true)
				}
				return nil
			},
		)
		if err == nil {
			c.Redirect(302, "/")
			c.Abort()
			return
		}
	}

	c.HTML(200, "admin_server.html", gin.H{
		"server": server,
		"error":  err,
	})
}

func handleAdminRestartServer(c *gin.Context) {
	ctl := getController(c)
	if ctl != nil {
		if err := ctl.Apply(true); err != nil {
			htmlError(c, err)
			return
		}
	}

	c.Redirect(302, "/admin")
}
