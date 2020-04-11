package service

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"xojoc.pw/useragent"
)

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
