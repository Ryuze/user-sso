package handler

import (
	"github.com/gin-gonic/gin"
	rest "github.com/ideal-tekno-solusi/sso/api/auth"
	"github.com/ideal-tekno-solusi/sso/bootstrap"
	"github.com/jackc/pgx/v5"
)

type RestService struct {
	db *pgx.Conn
}

func RestRegister(r *gin.Engine, cfg *bootstrap.Container) {
	rest.Router(r, &RestService{
		db: cfg.Db(),
	})
}
