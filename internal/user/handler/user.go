package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ideal-tekno-solusi/sso/api/user/operation"
	database "github.com/ideal-tekno-solusi/sso/database/main"
	"github.com/ideal-tekno-solusi/sso/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
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

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 32)
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
	queries := database.New(r.db)

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

	ctx.Status(http.StatusNoContent)
}
