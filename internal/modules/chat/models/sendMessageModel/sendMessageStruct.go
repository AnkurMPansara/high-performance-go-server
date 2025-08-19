package sendMessageModel

import "time"

type ApiInput struct {
	UserId  int `json:"user_id"`
	Message string `json:"message"`
}

type ApiData struct {
	StartTime time.Time
	UserId    int
	Message   string
	Reply     string
	Error     string
	Code      int
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Respose string `json:"response"`
	Error   string `json:"error"`
}
