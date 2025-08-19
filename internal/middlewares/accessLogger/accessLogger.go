package accessLogger

import (
	"backend-server/utilities/globalUtility"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func AccessLog(ginCtx *gin.Context) {
	logData := make(map[string]interface{})
	logData["IP"] = ginCtx.ClientIP()
	logData["PATH"] = ginCtx.FullPath()
	logData["TIMESTAMP"] = globalUtility.ConvertValueToString(time.Now())
	logData["CODE"] = ginCtx.Writer.Status()
	logData["METHOD"] = ginCtx.Request.Method
	logData["HEADERS"] = globalUtility.ConvertValueToString(ginCtx.Request.Header)
	if bodyBytes, err := io.ReadAll(ginCtx.Request.Body); err == nil {
		logData["BODY"] = string(bodyBytes)
		ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	globalUtility.CreateAccessLogs(logData)

	ginCtx.Next()
}