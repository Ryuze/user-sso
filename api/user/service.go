package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api/user/operation"
)

type Service interface {
	Register(ctx *gin.Context, params *operation.RegisterRequest)
	Login(ctx *gin.Context, params *operation.LoginRequest)
	Authorization(ctx *gin.Context, params *operation.AuthorizationRequest)
	Refresh(ctx *gin.Context, params *operation.RefreshRequest)
}
