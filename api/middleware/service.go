package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	VerifyToken(ctx *gin.Context, params *TokenRequest)
}
