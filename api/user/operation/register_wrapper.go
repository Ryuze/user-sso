//go:generate stringer -type=Gender -linecomment -output=gender_string.go
package operation

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/util"
)

type Gender int

const (
	UNIDENTIFIED Gender = iota // UNIDENTIFIED
	MALE                       // MALE
	FEMALE                     // FEMALE
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Dob      string `json:"dob"`
	Gender   Gender `json:"gender"`
}

func RegisterWrapper(handler func(ctx *gin.Context, params *RegisterRequest)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := RegisterRequest{}

		err := ctx.BindJSON(&params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		decryptPass, err := util.DecryptJwe(params.Password)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		params.Password = *decryptPass

		err = validateRegisterReq(params)
		if err != nil {
			util.SendProblemDetailJson(ctx, http.StatusInternalServerError, err.Error(), ctx.FullPath(), uuid.NewString())

			return
		}

		params.Password = *decryptPass

		handler(ctx, &params)

		ctx.Next()
	}
}

func validateRegisterReq(params RegisterRequest) error {
	if params.Username == "" {
		return errors.New("username can't be empty")
	}

	if len(params.Username) > 25 {
		return errors.New("username is too long")
	}

	if params.Password == "" {
		return errors.New("password can't be empty")
	}

	if len(params.Password) > 72 {
		return errors.New("password is to long")
	}

	if params.Name == "" {
		return errors.New("name can't be empty")
	}

	if len(params.Email) > 120 {
		return errors.New("name is too long")
	}

	if params.Email == "" {
		return errors.New("email can't be empty")
	}

	if len(params.Email) > 50 {
		return errors.New("email is too long")
	}

	if params.Dob == "" {
		return errors.New("dob can't be empty")
	}

	_, err := time.Parse("2006-01-02", params.Dob)
	if err != nil {
		return errors.New("dob not in correct format yyyy-mm-dd")
	}

	if params.Gender == 0 {
		return errors.New("gender should be choose")
	}

	if params.Gender > 2 {
		return errors.New("gender not valid")
	}

	return nil
}
