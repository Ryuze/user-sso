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

		response := util.GenerateProblemJson(http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusInternalServerError, response)

		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 32)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash password with error: %v", err)
		logrus.Warn(errorMessage)

		response := util.GenerateProblemJson(http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusInternalServerError, response)

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

		response := util.GenerateProblemJson(http.StatusInternalServerError, errorMessage, ctx.FullPath(), uuid.NewString())

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusInternalServerError, response)

		return
	}

	ctx.Status(http.StatusNoContent)
}
