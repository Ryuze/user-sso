package operation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/util"
)

type RefreshRequest struct {
	Jwt     string
	Service string `uri:"service"`
}

func RefreshWrapper(handler func(ctx *gin.Context, params *RefreshRequest)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := RefreshRequest{}

		err := ctx.BindUri(&params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		err = validateRefreshReq(params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusBadRequest, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		cookie, err := ctx.Cookie(fmt.Sprintf("jwt-%v", strings.ToLower(params.Service)))
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		params.Jwt = cookie
		params.Service = strings.ToLower(params.Service)

		handler(ctx, &params)

		ctx.Next()
	}
}

func validateRefreshReq(params RefreshRequest) error {
	if params.Service == "" {
		return errors.New("username can't be empty")
	}

	return nil
}
