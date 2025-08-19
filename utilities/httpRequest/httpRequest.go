package httpRequest

import (
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func MakeHttpRequest(ginCtx *gin.Context, apiService string, method string, payload interface{}, headers http.Header, timeout int) (apiCode int, apiResponse []byte, executiontime time.Duration, apiErr error) {
	start := time.Now()
	defer func() {
		if panic := recover(); panic != nil {
			apiCode = http.StatusInternalServerError
			apiResponse = nil
			apiErr = fmt.Errorf("some issue in makehttprequest: %v", panic)
		}
		executiontime = time.Since(start)
	}()
	apiUrl := configuration.GetConfigStringValue(apiService)
	httpCtx := ginCtx.Request.Context()
	var requestBody io.Reader
	if payload != nil && method != http.MethodGet && method != http.MethodHead {
		contentType := headers.Get("Content-Type")
		switch contentType {
		case "application/json":
			bodyBytes, bodyBytesErr := json.Marshal(payload)
			if bodyBytesErr != nil {
				apiCode = http.StatusInternalServerError
				apiResponse = nil
				apiErr = fmt.Errorf("some issue in marshaling payload in makehttprequest: %v", bodyBytesErr)
				return
			}
			requestBody = bytes.NewBuffer(bodyBytes)
		case "application/x-www-form-urlencoded":
			if bodyData, isBodyDataAvailable := payload.(map[string]string); isBodyDataAvailable {
				vals := url.Values{}
				for key, value := range bodyData {
					vals.Set(key, value)
				}
				requestBody = strings.NewReader(vals.Encode())
			} else {
				apiCode = http.StatusInternalServerError
				apiResponse = nil
				apiErr = fmt.Errorf("invalid payload for form encoding in makehttprequest")
				return
			}
		case "multipart/form-data":
			if bodyData, isBodyDataAvailable := payload.(map[string]interface{}); isBodyDataAvailable {
				var byteBuffer bytes.Buffer
				multiFormWriter := multipart.NewWriter(&byteBuffer)
				for key, value := range bodyData {
					switch v := value.(type) {
					case string:
						_ = multiFormWriter.WriteField(key, v)
					case []byte:
						part, _ := multiFormWriter.CreateFormFile(key, key)
						part.Write(v)
					case *os.File:
						part, _ := multiFormWriter.CreateFormFile(key, v.Name())
						io.Copy(part, v)
					case io.Reader:
						part, _ := multiFormWriter.CreateFormFile(key, key)
						io.Copy(part, v)
					default:
						_ = multiFormWriter.WriteField(key, globalUtility.ConvertValueToString(v))
					}
				}
				multiFormWriter.Close()
				requestBody = &byteBuffer
				headers.Set("Content-Type", multiFormWriter.FormDataContentType())
			} else {
				apiCode = http.StatusInternalServerError
				apiResponse = nil
				apiErr = fmt.Errorf("invalid payload for form-data in makehttprequest")
				return
			}
		case "text/plain":
			switch v := payload.(type) {
			case string:
				requestBody = strings.NewReader(v)
			case []byte:
				requestBody = bytes.NewBuffer(v)
			default:
				apiErr = fmt.Errorf("payload must be string or []byte for %v", contentType)
				return
			}
		case "application/octet-stream":
			switch v := payload.(type) {
			case string:
				requestBody = strings.NewReader(v)
			case []byte:
				requestBody = bytes.NewBuffer(v)
			default:
				apiErr = fmt.Errorf("payload must be string or []byte for %v", contentType)
				return
			}
		default:
			bodyBytes, bodyBytesErr := json.Marshal(payload)
			if bodyBytesErr != nil {
				apiCode = http.StatusInternalServerError
				apiResponse = nil
				apiErr = fmt.Errorf("some issue in marshaling payload in makehttprequest: %v", bodyBytesErr)
				return
			}
			requestBody = bytes.NewBuffer(bodyBytes)
		}
	}
	httpReq, httpReqErr := http.NewRequestWithContext(httpCtx, method, apiUrl, requestBody)
	if httpReqErr != nil {
		apiCode = http.StatusInternalServerError
		apiResponse = nil
		apiErr = fmt.Errorf("some issue in creating new request in makehttprequest: %v", httpReqErr)
		return
	}
	if headers != nil {
		httpReq.Header = headers
	}
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	httpReq.Close = true

	apiResult, apiResultErr := httpClient.Do(httpReq)
	if apiResultErr != nil {
		apiCode = http.StatusInternalServerError
		apiResponse = nil
		apiErr = fmt.Errorf("some issue while sending request in makehttprequest: %v", apiResultErr)
		return
	}

	defer apiResult.Body.Close()

	apiResp, apiRespErr := io.ReadAll(apiResult.Body)
	if apiRespErr != nil {
		apiCode = http.StatusInternalServerError
		apiResponse = nil
		apiErr = fmt.Errorf("some issue while reading response payload in makehttprequest: %v", apiRespErr)
		return
	}

	apiCode = apiResult.StatusCode
	apiResponse = apiResp
	apiErr = nil
	return
}
