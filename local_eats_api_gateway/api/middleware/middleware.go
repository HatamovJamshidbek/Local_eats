package middleware

import (
	t "api_get_way/api/token"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func MiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		url := ctx.Request.URL.Path
		fmt.Println("++++++++++", token)
		if strings.Contains(url, "swagger") || url == "/user/login" || url == "/user/register" {
			ctx.Next()
			return
		} else if _, err := t.ExtractClaim(token); err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.Next()

	}
}
