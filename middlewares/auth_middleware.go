package middlewares

import (
	"TestProject/helpers"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			helpers.ErrorResponse(c, 401, "Authorization header is required")
			c.Abort()
			return
		}
		claims, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			helpers.ErrorResponse(c, 401, "Invalid token")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
