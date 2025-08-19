package createTokenModel

import (
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	sessionId := strings.TrimRight(base64.URLEncoding.EncodeToString(bytes), "=")
	return sessionId, nil
}

func GenerateToken(apiData ApiData) (string, error) {
	tokenHeader := JwtTokenHeader{
		Algorithm : "HS256",
		Type: "JWT",
	}
	tokenPayload := JwtTokenPayLoad{
		Issuer: "backend-server",
		Subject: globalUtility.ConvertValueToString(apiData.UserId),
		Audience: "user",
		ExpirationTime: int(apiData.ExpirationTime.Unix()),
		NotBeforeTime: int(apiData.ValidationStartTime.Unix()),
		IssuedAt: int(time.Now().Unix()),
		JwtId: apiData.SessionId,
	}
	var jwtHeader string
	var jwtPayload string
	tokenHeaderJson, tokenHeaderJsonErr := json.Marshal(tokenHeader)
	if tokenHeaderJsonErr != nil {
		return "", tokenHeaderJsonErr
	}
	jwtHeader = base64.RawURLEncoding.EncodeToString(tokenHeaderJson)
	tokenPayloadJson, tokenPayloadJsonErr := json.Marshal(tokenPayload)
	if tokenPayloadJsonErr != nil {
		return "", tokenPayloadJsonErr
	}
	jwtPayload = base64.RawURLEncoding.EncodeToString(tokenPayloadJson)
	secretKey := configuration.GetConfigStringValue("authentication_secret_key_user")
	signInput := jwtHeader + "." + jwtPayload
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(signInput))
	signature := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	return jwtHeader + "." + jwtPayload + "." + signature, nil
}

func CreateLogs(ginCtx *gin.Context, apiInput ApiInput, apiData ApiData) {
	logData := make(map[string]interface{})
	logData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_create_token")
	logData["TIMESTAMP"] = globalUtility.ConvertValueToString(apiData.StartTime)
	logData["EXECUTION_TIME"] = time.Since(apiData.StartTime).Microseconds()
	logData["CODE"] = apiData.Code
	logData["USER_ID"] = apiData.UserId
	logData["ERROR"] = apiData.Error
	logData["GENERATED_TOKEN"] = apiData.GeneratedToken
	globalUtility.CreateApplicationLogs(logData)
}
