package customRecovery

import (
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func HandlePanic(ginCtx *gin.Context) {
	defer func() {
		if panicRecovery := recover(); panicRecovery != nil {
			stackTrace := string(debug.Stack())
			errorLogData := make(map[string]interface{})
			errorLogData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_custom_recovery")
			errorLogData["CODE"] = http.StatusInternalServerError
			errorLogData["ERROR"] = globalUtility.ConvertValueToString(panicRecovery)
			errorLogData["STACK"] = stackTrace
			ginCtx.Abort()
			globalUtility.CreateApplicationLogs(errorLogData)
		}
	}()
	ginCtx.Next()
}