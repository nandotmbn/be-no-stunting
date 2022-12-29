package helpers

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(id string, role string) (tokenString string, err error) {
	claims := &JWTClaim{
		Id:             id,
		Role:           role,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (val string) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})

	if token == nil {
		return "Still error"
	}

	for _, val := range claims {
		return fmt.Sprintf("%s", val)
	}
	return err.Error()
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
