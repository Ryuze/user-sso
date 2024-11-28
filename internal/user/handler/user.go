package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/api/user/operation"
	database "github.com/ideal-tekno-solusi/sso/database/main"
	"github.com/ideal-tekno-solusi/sso/internal/user/entity/response"
	"github.com/ideal-tekno-solusi/sso/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (r *RestService) Register(ctx *gin.Context, params *operation.RegisterRequest) {
	queries := database.New(r.db)

	user := database.CreateUserReturnIdParams{
		Username: params.Username,
		Name:     params.Name,
		Email:    params.Email,
		Dob:      params.Dob,
		Gender:   params.Gender.String(),
	}

	insertedUserId, err := queries.CreateUserReturnId(ctx, user)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to insert new user with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 16)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash password with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	password := database.CreatePasswordReturnUserIdParams{
		UserID: pgtype.Int4{
			Int32: insertedUserId,
			Valid: true,
		},
		Password: string(hashPassword),
	}

	_, err = queries.CreatePasswordReturnUserId(ctx, password)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to insert password to new user with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	ctx.Status(http.StatusNoContent)
}

func (r *RestService) Login(ctx *gin.Context, params *operation.LoginRequest) {
	domain := viper.GetString("config.domain")

	queries := database.New(r.db)

	res := response.LoginResponse{}

	password, err := queries.GetUserLatestPassword(ctx, params.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to fetch user password list with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(params.Password))
	if err != nil {
		errorMessage := "password not match"
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	user, err := queries.GetUser(ctx, params.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get user information with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	services := strings.Split(user.AllowedServices.String, ",")
	if !slices.Contains(services, params.Service) {
		errorMessage := fmt.Sprintf("service %v is not authorized for current user", params.Service)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	token, time, err := util.BuildUserJwt(user)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to build token with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	refresh, _, err := util.BuildRefreshJwt(user.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to build token with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	sign, err := util.SignJwt(*token, params.Service)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to sign token with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	refreshSign, err := util.SignJwt(*refresh, params.Service)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to sign token with error: %v", err)
		logrus.Warn(errorMessage)

		util.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		return
	}

	res.Authorization = fmt.Sprintf("Bearer %v", *sign)
	res.Time = int(time.Seconds())

	ctx.SetCookie(fmt.Sprintf("jwt-%v", params.Service), *refreshSign, 60*60*24, "/", domain, false, true)

	ctx.JSON(http.StatusOK, res)
}
