package getGreetingsModel

import (
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"time"

	"github.com/gin-gonic/gin"
)

func FetchGreetings(input string) string {
	return "Hello User! I heard you said " + input
}

func CreateLogs(ginCtx *gin.Context, apiInput ApiInput, apiData ApiData) {
	logData := make(map[string]interface{})
	logData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_get_greetings")
	logData["TIMESTAMP"] = globalUtility.ConvertValueToString(apiData.StartTime)
	logData["EXECUTION_TIME"] = time.Since(apiData.StartTime).Microseconds()
	logData["CODE"] = apiData.Code
	logData["REPLY"] = apiData.Reply
	logData["ERROR"] = apiData.Error
	logData["CHAT_RESPONSE"] = apiData.ChatResponse
	globalUtility.CreateApplicationLogs(logData)
}