package httpx

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRealIP(c *gin.Context) string {
	var realIP string
	forwardIP := c.Request.Header.Get("X-Forwarded-For")
	if forwardIP == "" {
		realIP = c.Request.Header.Get("X-Real-Ip")
	} else {
		realIP = strings.Split(forwardIP, ",")[0]
	}
	if realIP != "" {
		return realIP
	}
	return c.ClientIP()
}
