package service

import (
	"fmt"
	"log"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"golang.org/x/oauth2"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"xojoc.pw/useragent"

	"github.com/sosedoff/x/wireguard-manager/generate"
	"github.com/sosedoff/x/wireguard-manager/middleware"
	"github.com/sosedoff/x/wireguard-manager/model"
)

func New(config *Config) (*gin.Engine, error) {
	router := gin.New()

	tpl, err := loadTemplate(config.AssetFS)
	if err != nil {
		return nil, err
	}
	router.SetHTMLTemplate(tpl)

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(middleware.RejectBots())
	router.Use(middleware.CookieSession(config.CookieName, config.CookieSecret))

	router.Use(setAuth(config))
	router.Use(setStore(config))
	router.Use(setController(config))

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
	admin.GET("/peers", handleAdminPeers)

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
		c.String(400, err.Error())
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		badRequest(c, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		badRequest(c, "Request failed")
		return
	}

	info, err := googleAuthFromResponse(resp)
	if err != nil {
		badRequest(c, err)
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
			badRequest(c, "Email address is not permitted")
			return
		}
	}

	store := getStore(c)
	user, err := store.FindUserByEmail(info.Email)
	if err != nil {
		badRequest(c, err)
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
			badRequest(c, err)
			return
		}
		if count == 0 {
			user.Role = roleAdmin
		}

		if err := store.CreateUser(user); err != nil {
			badRequest(c, err)
			return
		}
	}

	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()

	c.Redirect(302, "/")
}

func handleSignout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("uid")
	session.Save()

	c.Redirect(302, "/")
}

func handleIndex(c *gin.Context) {
	store := getStore(c)
	user := getUser(c)

	devices, err := store.GetDevicesByUser(user.ID)
	if err != nil {
		badRequest(c, err)
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

	server, err := store.FindServer()
	if err != nil {
		htmlError(c, err)
		return
	}
	if server == nil {
		htmlError(c, "No server found!")
		return
	}

	ipfirst, iplast, err := server.IPV4Range()
	if err != nil {
		htmlError(c, err)
		return
	}
	ipfirst = cidr.Inc(ipfirst)

	allocated, err := store.AllocatedIPV4()
	if err != nil {
		htmlError(c, err)
		return
	}

	cur := ipfirst

	for {
		taken := false
		for _, val := range allocated {
			if val == cur.String() {
				taken = true
				break
			}
		}

		if !taken {
			break
		}

		if cur.Equal(iplast) {
			log.Fatal("not more addrs")
		}

		cur = cidr.Inc(cur)
	}

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		htmlError(c, err)
		return
	}

	device := &model.Device{
		UserID:              user.ID,
		Name:                c.Request.FormValue("name"),
		OS:                  c.Request.FormValue("os"),
		Enabled:             true,
		PrivateKey:          key.String(),
		PublicKey:           key.PublicKey().String(),
		IPV4:                cur.String(),
		PersistentKeepalive: 60,
	}

	if err := store.CreateDevice(device); err != nil {
		htmlError(c, err)
		return
	}

	getController(c).ScheduleApply()

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

	device, err := store.FindUserDevice(user, c.Param("id"))
	if err != nil {
		badRequest(c, err)
		return
	}
	if device == nil {
		c.Redirect(302, "/")
		return
	}

	if err := store.DeleteUserDevice(user, device); err != nil {
		badRequest(c, err)
		return
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

	config, err := generate.ServerConfig(store)
	if err != nil {
		htmlError(c, err)
		return
	}

	c.HTML(200, "admin_index.html", gin.H{
		"devices": devices,
		"users":   users,
		"server":  server,
		"config":  string(config),
	})
}

func handleAdminPeers(c *gin.Context) {
	server, err := getStore(c).FindServer()
	if err != nil {
		htmlError(c, err)
		return
	}

	client, err := wgctrl.New()
	if err != nil {
		htmlError(c, err)
		return
	}

	device, err := client.Device(server.Interface)
	if err != nil {
		htmlError(c, err)
		return
	}

	peerInfos := make([]map[string]interface{}, len(device.Peers))

	for idx, p := range device.Peers {
		ips := make([]string, len(p.AllowedIPs))
		for ipidx, ip := range p.AllowedIPs {
			ips[ipidx] = ip.String()
		}

		peerInfos[idx] = map[string]interface{}{
			"PublicKey":         p.PublicKey.String(),
			"Endpoint":          p.Endpoint.String(),
			"LastHandshakeTime": p.LastHandshakeTime,
			"ReceiveBytes":      p.ReceiveBytes,
			"TransmitBytes":     p.TransmitBytes,
			"AllowedIPs":        ips,
		}
	}

	c.HTML(200, "admin_peers.html", gin.H{"peers": peerInfos})
}
