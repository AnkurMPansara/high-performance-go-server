package sendMessageModel

import (
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateLogs(ginCtx *gin.Context, apiInput ApiInput, apiData ApiData) {
	logData := make(map[string]interface{})
	logData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_create_token")
	logData["TIMESTAMP"] = globalUtility.ConvertValueToString(apiData.StartTime)
	logData["EXECUTION_TIME"] = time.Since(apiData.StartTime).Microseconds()
	logData["CODE"] = apiData.Code
	logData["MESSAGE"] = apiData.Message
	logData["ERROR"] = apiData.Error
	logData["REPLY"] = apiData.Reply
	logData["USER_ID"] = apiData.UserId
	globalUtility.CreateApplicationLogs(logData)
}