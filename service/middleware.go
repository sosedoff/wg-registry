package service

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"xojoc.pw/useragent"
)

func forceHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS != nil {
			c.Next()
			return
		}

		destination := &url.URL{
			Scheme:   "https",
			Host:     c.Request.Host,
			Path:     c.Request.URL.Path,
			RawQuery: c.Request.URL.RawQuery,
		}

		c.Redirect(301, destination.String())
		c.Abort()
	}
}

func rejectBots() gin.HandlerFunc {
	return func(c *gin.Context) {
		agent := useragent.Parse(c.Request.UserAgent())
		if agent == nil {
			c.Next()
			return
		}

		if agent.Type != useragent.Browser {
			c.AbortWithError(http.StatusForbidden, errors.New("forbidden"))
		}
	}
}

func cookieSession(name, secret string) gin.HandlerFunc {
	return sessions.Sessions(name, cookie.NewStore([]byte(secret)))
}
