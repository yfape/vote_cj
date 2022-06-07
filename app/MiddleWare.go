package app

import (
	"net/http"
	"vote_cj/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		//vcalidate token formate
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 460, "msg": "非法请求"})
			ctx.Abort()
			return
		}

		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 461, "msg": "非法请求"})
			ctx.Abort()
			return
		}

		checkcode := claims.CheckCode
		if CreateCheckCode(claims.OpenId) != checkcode {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 462, "msg": "非法请求"})
			ctx.Abort()
			return
		}

		ctx.Set("openid", claims.OpenId)
		ctx.Next()
	}
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.System.Secret), nil
	})
	return token, Claims, err
}
