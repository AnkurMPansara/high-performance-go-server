package authentication

import (
	"backend-server/utilities/configuration"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JwtTokenHeader struct {
	Algorithm string `json:"alg"`
	Type string `json:"typ"`
}

type JwtTokenPayLoad struct {
	Issuer         string `json:"iss"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	ExpirationTime int    `json:"exp"`
	NotBeforeTime  int    `json:"nbf"`
	IssuedAt       int    `json:"iat"`
	JwtId          string `json:"jti"`
}

func AuthenticateRequest(ginCtx *gin.Context) {
	token := ginCtx.GetHeader("Authorization")
	userId := ginCtx.GetHeader("UserId")
	if token == "" {
		ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
			"code": http.StatusUnauthorized,
			"status": "Failure",
			"response": "",
			"error": "Unauthorized Request",
		})
	}
	if !validateToken(token, userId){
		ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
			"code": http.StatusUnauthorized,
			"status": "Failure",
			"response": "",
			"error": "Invalid Authorization",
		})
	}
	ginCtx.Next()
}

func validateToken(token string, userId string) bool {
	token = strings.ReplaceAll(token, "Bearer ", "")
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 3 {
		return false
	}
	headerBytes, headerDecodeErr := base64.RawURLEncoding.DecodeString(tokenParts[0])
	if headerDecodeErr != nil {
		return false
	}
	var tokenHeader JwtTokenHeader
	tokenHeaderErr := json.Unmarshal(headerBytes, &tokenHeader)
	if tokenHeaderErr != nil {
		return false
	}

	payloadBytes, payloadDecodeErr := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if payloadDecodeErr != nil {
		return false
	}
	var tokenPayload JwtTokenPayLoad
	tokenPayloadErr := json.Unmarshal(payloadBytes, &tokenPayload)
	if tokenPayloadErr != nil {
		return false
	}

	if tokenHeader.Type == "JWT" {
		if tokenPayload.NotBeforeTime > int(time.Now().Unix()) {
			return false
		}
		if tokenPayload.ExpirationTime < int(time.Now().Unix()) {
			return false
		}
		if tokenPayload.Audience == "user" {
			if tokenPayload.Subject != userId {
				return false
			}
			if tokenHeader.Algorithm == "HS256" {
				signInput := tokenParts[0] + "." + tokenParts[1]
				secretKey := configuration.GetConfigStringValue("authentication_secret_key_user")
				hash := hmac.New(sha256.New, []byte(secretKey))
				hash.Write([]byte(signInput))
				signature := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
				if tokenParts[2] != signature {
					return false
				}
			} else {
				return false
			}
		} else if tokenPayload.Audience == "server" {
			if tokenHeader.Algorithm == "HS256" {
				signInput := tokenParts[0] + "." + tokenParts[1]
				secretKey := configuration.GetConfigStringValue("authentication_secret_key_server")
				hash := hmac.New(sha256.New, []byte(secretKey))
				hash.Write([]byte(signInput))
				signature := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
				if tokenParts[2] != signature {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
	
	return true
}