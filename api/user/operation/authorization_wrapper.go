package operation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/util"
)

type AuthorizationRequest struct {
	Service string `json:"service"`
	Token   string `json:"token"`
}

func AuthorizationWrapper(handler func(ctx *gin.Context, params *AuthorizationRequest)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := AuthorizationRequest{}

		err := ctx.BindJSON(&params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		err = validateAuthorizationReq(params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusBadRequest, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		handler(ctx, &params)

		ctx.Next()
	}
}

func validateAuthorizationReq(params AuthorizationRequest) error {
	if params.Service == "" {
		return errors.New("service can't be empty")
	}

	if len(params.Service) > 50 {
		return errors.New("service is too long")
	}

	if params.Token == "" {
		return errors.New("token can't be empty")
	}

	return nil
}
