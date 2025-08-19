package sendMessageModel

import (
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"backend-server/utilities/httpRequest"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sendMessageToLLM(ginCtx *gin.Context, apiData *ApiData) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("Issue in sendMessageToLLM: %v", panicErr)
			return
		}
	}()
	llmApiKey := configuration.GetConfigStringValue("API_KEY_LLM")
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer " + llmApiKey)
	payload := map[string]interface{}{
        "model": "mistral-large-latest",
        "messages": []map[string]interface{}{
            {
                "role":    "system",
                "content": "You are a helpful assistant.",
            },
            {
                "role":    "user",
                "content": apiData.Message,
            },
        },
        "temperature": 0.7,
        "top_p":       1.0,
    }
	apiCode, apiResponse, executionTime, apiErr := httpRequest.MakeHttpRequest(ginCtx, "llm_chat_api", "POST", payload, headers, 60000)
	if apiErr != nil {
		err = apiErr
		return
	}
	apiData.ExecTimes["ext_exec1"] = executionTime
	if apiCode == http.StatusOK {
		var aiResponse ChatCompletionResponse
		unmarshalErr := json.Unmarshal(apiResponse, &aiResponse)
		if unmarshalErr != nil {
			err = fmt.Errorf("some issue in unmarshaling api response: %v", unmarshalErr)
			return
		}
		apiData.Reply = aiResponse.Choices[0].Message.Content
	}
	return nil
}

func CreateLogs(ginCtx *gin.Context, apiInput ApiInput, apiData ApiData) {
	logData := make(map[string]interface{})
	logData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_create_token")
	logData["TIMESTAMP"] = globalUtility.ConvertValueToString(apiData.StartTime)
	logData["EXECUTION_TIME"] = time.Since(apiData.StartTime).Microseconds()
	if len(apiData.ExecTimes) != 0 {
		for key, val := range apiData.ExecTimes {
			logData[key] = globalUtility.ConvertValueToInt(val)
		}
	}
	logData["CODE"] = apiData.Code
	logData["MESSAGE"] = apiData.Message
	logData["ERROR"] = apiData.Error
	logData["REPLY"] = apiData.Reply
	logData["USER_ID"] = apiData.UserId
	globalUtility.CreateApplicationLogs(logData)
}