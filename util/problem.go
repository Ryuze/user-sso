package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Problem struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   interface{} `json:"detail"`
	Instance string      `json:"instance"`
	Guid     string      `json:"guid"`
}

func GenerateProblemJson(statusCode int, message, instance, guid string) Problem {
	return Problem{
		Title:    http.StatusText(statusCode),
		Status:   statusCode,
		Detail:   message,
		Instance: instance,
		Guid:     guid,
	}
}

func SendProblemDetailJson(ctx *gin.Context, statusCode int, message, instance, guid string) {
	problem := GenerateProblemJson(statusCode, message, instance, guid)

	ctx.Header("Content-Type", "application/problem+json")
	ctx.JSON(statusCode, problem)
}
