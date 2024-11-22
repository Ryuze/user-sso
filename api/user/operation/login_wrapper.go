package operation

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/util"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Service  string `json:"service"`
}

func LoginWrapper(handler func(ctx *gin.Context, params *LoginRequest)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := LoginRequest{}

		err := ctx.BindJSON(&params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		err = validateLoginReq(params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusBadRequest, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		decryptPass, err := util.DecryptJwe(params.Password, strings.ToLower(params.Service))
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		params.Password = *decryptPass

		handler(ctx, &params)

		ctx.Next()
	}
}

func validateLoginReq(params LoginRequest) error {
	if params.Username == "" {
		return errors.New("username can't be empty")
	}

	if len(params.Username) > 25 {
		return errors.New("username is too long")
	}

	if params.Password == "" {
		return errors.New("password can't be empty")
	}

	if params.Service == "" {
		return errors.New("service can't be empty")
	}

	return nil
}
