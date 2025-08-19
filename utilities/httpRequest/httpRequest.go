package httpRequest

import (
	"backend-server/utilities/configuration"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MakeHttpRequest(ginCtx *gin.Context,apiService string, method string, payload interface{}, headers http.Header, timeout int) (apiCode int, apiResponse []byte, executiontime time.Duration, apiErr error) {
	start := time.Now()
	defer func() {
		if panic := recover(); panic != nil {
			apiCode = http.StatusInternalServerError
			apiResponse = nil
			apiErr = fmt.Errorf("Some issue in MakeHttpRequest: %s", panic)
		}
		executiontime = time.Since(start)
	}()
	url := configuration.GetConfigStringValue(apiService)
	httpCtx := ginCtx.Request.Context()
	requestBody, requestBodyErr := json.Marshal(payload)
	if requestBodyErr != nil {
		return
	}
	httpReq, httpReqErr := http.NewRequestWithContext(httpCtx, method, url, bytes.NewBuffer(requestBody))
	if httpReqErr != nil {

	}
	return
}