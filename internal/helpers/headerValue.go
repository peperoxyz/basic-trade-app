package helpers

import "github.com/gin-gonic/gin"

// to get the content-type dari header di request API client

func GetContentType(ctx *gin.Context) string {
	return ctx.Request.Header.Get("Content-Type")
}