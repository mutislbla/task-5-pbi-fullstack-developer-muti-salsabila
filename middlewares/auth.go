package middlewares

import (
	"net/http"
	"strings"
	"task5-pbi/app"
	"task5-pbi/database"
	"task5-pbi/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenClaims struct {
	ID int
	jwt.StandardClaims
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "No Token Available"})
		}
		accessToken := strings.Split(authHeader, " ")[1]
		tokenClaims := helpers.ParseToken(accessToken, c)
		var user app.User
		if err := database.DB.Where("ID=?", tokenClaims.ID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized Access"})
		}
		c.Set("user", user)
		c.Next()
	}
}
