package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/api/user/operation"
	database "github.com/ideal-tekno-solusi/sso/database/main"
	"github.com/ideal-tekno-solusi/sso/internal/user/entity/response"
	"github.com/ideal-tekno-solusi/sso/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func (r *RestService) Authorization(ctx *gin.Context, params *operation.AuthorizationRequest) {
	queries := database.New(r.db)

	res := response.Authorization{}

	service, err := queries.GetService(ctx, params.Service)
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
		ctx.JSON(http.StatusOK, res)
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

	res.Approved = true

	ctx.JSON(http.StatusOK, res)
}

func (r *RestService) Refresh(ctx *gin.Context, params *operation.RefreshRequest) {
	ctx.JSON(http.StatusOK, params)
}
