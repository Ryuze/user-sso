package util

import "net/http"

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
