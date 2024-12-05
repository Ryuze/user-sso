package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/bootstrap"
	database "github.com/ideal-tekno-solusi/sso/database/main"
	"github.com/ideal-tekno-solusi/sso/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type TokenRequest struct {
	Token string `header:"token"`
}

// ? verify token ini di hardcode khusus untuk verify service user
func VerifyToken(cfg *bootstrap.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := TokenRequest{}
		db := cfg.Db()
		queries := database.New(db)

		err := ctx.BindHeader(&params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		service, err := queries.GetService(ctx, "user")
		if err != nil {
			errorMessage := fmt.Sprintf("failed to get service with error: %v", err)
			logrus.Warn(errorMessage)

			switch err {
			case pgx.ErrNoRows:
				util.SendProblemDetailJson(ctx, http.StatusNotFound, errorMessage, ctx.FullPath(), uuid.NewString())
				return
			default:
				util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())
				return
			}
		}
		if !service.Status {
			util.SendProblemDetailJson(ctx, http.StatusServiceUnavailable, "current service status is offline", ctx.FullPath(), uuid.NewString())
			return
		}

		token := strings.Split(params.Token, " ")
		switch token[0] {
		case "Bearer":
			_, err := util.VerifyJwt(token[1], service.PublicKey)
			if err != nil {
				errorMessage := "failed to validate token"
				logrus.Warn(errorMessage)

				util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

				return
			}
		default:
			errorMessage := fmt.Sprintf("token %v is not an accepted token", token[0])
			logrus.Warn(errorMessage)

			util.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.FullPath(), uuid.NewString())

			return
		}

		ctx.Next()
	}
}
