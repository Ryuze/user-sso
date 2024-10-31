package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api/user/operation"
)

func Router(r *gin.Engine, s Service) {
	r.POST("/v1/register", operation.RegisterWrapper(s.Register))
	r.POST("/v1/login", operation.LoginWrapper(s.Login))
}
