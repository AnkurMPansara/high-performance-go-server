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
	ExecTimes map[string]time.Duration
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Respose string `json:"response"`
	Error   string `json:"error"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Created int64  `json:"created"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Index int `json:"index"`
		Message struct {
			Content   string `json:"content"`
			Prefix    bool   `json:"prefix"`
			Role      string `json:"role"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments"`
				} `json:"function"`
				Index int `json:"index"`
			} `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}