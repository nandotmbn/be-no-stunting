package helpers

import (
	"be-no-stunting-v2/configs"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(configs.EnvSuperSecret())

type JWTClaim struct {
	Id   string    `json:"id"`
	Role string    `json:"role"`
	Time time.Time `json:"time"`
	jwt.StandardClaims
}

func GenerateJWT(id string, role string) (tokenString string, err error) {
	claims := &JWTClaim{
		Id:             id,
		Role:           role,
		Time:           time.Now(),
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (val string, err error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if token == nil {
		val = "Token is invalid"
		return
	}

	val = fmt.Sprintf("%s", claims["id"])
	return val, nil
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
