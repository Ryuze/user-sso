package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api/test"
	"github.com/ideal-tekno-solusi/sso/bootstrap"
	auth "github.com/ideal-tekno-solusi/sso/internal/auth/handler"
	user "github.com/ideal-tekno-solusi/sso/internal/user/handler"
)

func RegisterApi(r *gin.Engine, cfg *bootstrap.Container) {
	user.RestRegister(r, cfg)
	auth.RestRegister(r, cfg)
	test.Router(r)
}
