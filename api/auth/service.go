package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api/auth/operation"
)

type Service interface {
	Authorization(ctx *gin.Context, params *operation.AuthorizationRequest)
	Refresh(ctx *gin.Context, params *operation.RefreshRequest)
}
