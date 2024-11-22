package operation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/util"
)

type RefreshRequest struct {
	Jwt string
}

func RefreshWrapper(handler func(ctx *gin.Context, params *RefreshRequest)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := RefreshRequest{}

		cookie, err := ctx.Cookie("jwt")
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		params.Jwt = cookie

		handler(ctx, &params)

		ctx.Next()
	}
}
