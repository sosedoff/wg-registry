package service

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/sosedoff/wg-registry/controller"
	"github.com/sosedoff/wg-registry/model"
	"github.com/sosedoff/wg-registry/store"
)

var (
	roleUser  = "user"
	roleAdmin = "admin"

	ctxStoreKey      = "store"
	ctxUserKey       = "user"
	ctxDeviceKey     = "device"
	ctxControllerKey = "controller"
	ctxAuthKey       = "oauth"
)

func getSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

func getStore(c *gin.Context) store.Store {
	return c.MustGet(ctxStoreKey).(store.Store)
}

func getUser(c *gin.Context) *model.User {
	return c.MustGet(ctxUserKey).(*model.User)
}

func getDevice(c *gin.Context) *model.Device {
	return c.MustGet(ctxDeviceKey).(*model.Device)
}

func getController(c *gin.Context) *controller.Controller {
	return c.MustGet(ctxControllerKey).(*controller.Controller)
}

func getAuth(c *gin.Context) *AuthConfig {
	return c.MustGet(ctxAuthKey).(*AuthConfig)
}

func setStore(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ctxStoreKey, config.Store)
	}
}

func setController(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ctxControllerKey, config.Controller)
	}
}

func setAuth(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ctxAuthKey, &AuthConfig{
			domain:    config.ClientDomain,
			whitelist: config.ClientWhitelist,
			oauth: &oauth2.Config{
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				Scopes:       []string{"email", "profile"},
				Endpoint:     google.Endpoint,
			},
		})
	}
}

func requireUser(c *gin.Context) {
	session := getSession(c)
	store := getStore(c)

	uid := session.Get("uid")
	if uid == nil || uid == "" {
		c.Abort()
		c.Redirect(302, authStartPath)
		return
	}

	user, err := store.FindUserByID(uid)
	if err != nil {
		badRequest(c, err)
		return
	}
	if user == nil {
		c.Abort()
		c.Redirect(302, authStartPath)
		return
	}

	c.Set(ctxUserKey, user)
}

func requireAdmin(c *gin.Context) {
	user := getUser(c)

	if user.Role != roleAdmin {
		c.AbortWithError(403, errors.New("Permission denied"))
		return
	}
}

func requireDevice(c *gin.Context) {
	store := getStore(c)
	user := getUser(c)

	device, err := store.FindUserDevice(user, c.Param("id"))
	if err != nil {
		htmlError(c, err)
		return
	}

	if device == nil {
		htmlError(c, errors.New("Device does not exist"))
		return
	}

	c.Set(ctxDeviceKey, device)
}

func serveStaticAsset(fs *assets.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		if filepath.Ext(c.Request.URL.Path) == ".html" {
			c.AbortWithStatus(404)
			return
		}

		asset, ok := fs.Files[c.Request.URL.Path]
		if !ok {
			c.AbortWithStatus(404)
			return
		}
		c.Data(200, "text/plain", asset.Data)
	}
}

func getRequestScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func makeRedirectURL(req *http.Request, path string) string {
	host := strings.ToLower(req.Host)
	scheme := getRequestScheme(req)

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}
