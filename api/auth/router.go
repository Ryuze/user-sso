package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api/auth/operation"
)

func Router(r *gin.Engine, s Service) {
	r.POST("/v1/auth/authorization", operation.AuthorizationWrapper(s.Authorization))
	r.GET("/v1/auth/refresh/:service", operation.RefreshWrapper(s.Refresh))
}
