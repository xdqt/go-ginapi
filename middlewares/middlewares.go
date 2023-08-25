package middlewares

import (
	"net/http"

	"ginapi/utils/token"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, newtoken := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		// 设置表头
		c.Writer.Header().Set("Authorization", "JWT "+newtoken)
		c.SetCookie("Authorization", newtoken, 3600, "/", "localhost", true, true)
		c.Next()

	}
}
