package getGreetingsModel

import "time"

type ApiInput struct {
	Reply string `json:"reply"`
}

type ApiData struct {
	StartTime    time.Time
	Reply        string
	ChatResponse string
	Error        string
	Code         int
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Respose string `json:"response"`
	Error   string `json:"error"`
}