package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api/user/operation"
)

func Router(r *gin.Engine, s Service) {
	r.POST("/v1/user/register", operation.RegisterWrapper(s.Register))
	r.POST("/v1/user/login", operation.LoginWrapper(s.Login))
	r.POST("/v1/auth/authorization", operation.AuthorizationWrapper(s.Authorization))
	r.GET("/v1/auth/refresh", operation.RefreshWrapper(s.Refresh))
}
