package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"xojoc.pw/useragent"
)

var (
	errForbidden = errors.New("forbidden")
)

// RejectBots returns user agent handling middleware
func RejectBots() gin.HandlerFunc {
	return func(c *gin.Context) {
		agent := useragent.Parse(c.Request.UserAgent())
		if agent == nil {
			c.Next()
			return
		}

		if agent.Type != useragent.Browser {
			c.AbortWithError(http.StatusForbidden, errForbidden)
		}
	}
}

// CookieSession returns cookie handling middleware
func CookieSession(name, secret string) gin.HandlerFunc {
	return sessions.Sessions(name, cookie.NewStore([]byte(secret)))
}
