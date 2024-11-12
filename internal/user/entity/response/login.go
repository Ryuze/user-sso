package response

type LoginResponse struct {
	Authorization string `json:"authorization"`
	Time          int    `json:"time"`
}
