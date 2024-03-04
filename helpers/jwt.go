package helpers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	ID int
	jwt.StandardClaims
}

var secretKey = os.Getenv("SECRET_KEY")

func GenerateToken(id int, c *gin.Context) string {
	claims := Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
	}
	return tokenString
}

func ParseToken(signedToken string, c *gin.Context) *Claims {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		panic(err.Error())
	}

	claims, ok := token.Claims.(*Claims)
	fmt.Println(">>>", claims, ok)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse token"})
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token expired"})
	}

	return claims
}
