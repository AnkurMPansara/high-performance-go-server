package createTokenController

import (
	"backend-server/internal/modules/service/models/createTokenModel"
	"backend-server/utilities/globalUtility"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CreateTokenController struct {
	apiData     createTokenModel.ApiData
	apiInput    createTokenModel.ApiInput
	apiResponse createTokenModel.ApiResponse
}

func CreateToken(ginCtx *gin.Context) {
	controller := &CreateTokenController{}
	controller.apiData.StartTime = time.Now()
	controller.apiData.Code = http.StatusOK
	controller.perform(ginCtx)
}

func (controller *CreateTokenController) perform(ginCtx *gin.Context) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			errorMsg := "Issue in createToken API: " + fmt.Sprintf("%v", panicErr)
			controller.apiData.Error = errorMsg
			controller.apiData.Code = http.StatusInternalServerError
			controller.returnApiResp(ginCtx)
		}
		createTokenModel.CreateLogs(ginCtx, controller.apiInput, controller.apiData)
	}()
	if bindingErr := controller.bindInput(ginCtx); bindingErr != nil {
		controller.apiResponse.Code = http.StatusBadRequest
		controller.apiResponse.Status = "Failure"
		controller.apiResponse.Error = "Invalid Request Parameters"
		controller.apiData.Error = globalUtility.ConvertValueToString(bindingErr)
		controller.apiData.Code = http.StatusBadRequest
		controller.returnApiResp(ginCtx)
		return
	}
	if validationErr := controller.validateInput(); validationErr != nil {
		controller.apiResponse.Code = http.StatusBadRequest
		controller.apiResponse.Status = "Failure"
		controller.apiResponse.Error = globalUtility.ConvertValueToString(validationErr)
		controller.apiData.Error = globalUtility.ConvertValueToString(validationErr)
		controller.apiData.Code = http.StatusBadRequest
		controller.returnApiResp(ginCtx)
		return
	}
	sessionId, generateSessionErr := createTokenModel.GenerateSessionId()
	if generateSessionErr != nil {
		controller.apiResponse.Code = http.StatusInternalServerError
		controller.apiResponse.Status = "Failure"
		controller.apiResponse.Error = "Internal Server Error"
		controller.apiData.Error = "Some issue in generating sessionID: " + globalUtility.ConvertValueToString(generateSessionErr)
		controller.apiData.Code = http.StatusInternalServerError
		controller.returnApiResp(ginCtx)
		return
	}
	controller.apiData.SessionId = sessionId
	generatedToken, generatedTokenErr := createTokenModel.GenerateToken(controller.apiData)
	if generatedTokenErr != nil {
		controller.apiResponse.Code = http.StatusInternalServerError
		controller.apiResponse.Status = "Failure"
		controller.apiResponse.Error = "Internal Server Error"
		controller.apiData.Error = "Some issue in generating token: " + globalUtility.ConvertValueToString(generatedTokenErr)
		controller.apiData.Code = http.StatusInternalServerError
		controller.returnApiResp(ginCtx)
		return
	}
	controller.apiData.GeneratedToken = generatedToken

	controller.apiResponse.Code = http.StatusOK
	controller.apiResponse.Status = "Success"
	controller.apiResponse.Respose = createTokenModel.CreateTokenRespose{
		Token:  controller.apiData.GeneratedToken,
		UserId: controller.apiData.UserId,
	}
	controller.returnApiResp(ginCtx)
}

func (controller *CreateTokenController) returnApiResp(ginCtx *gin.Context) {
	if !ginCtx.Writer.Written() {
		ginCtx.JSON(controller.apiResponse.Code, controller.apiResponse)
	}
}

func (controller *CreateTokenController) bindInput(ginCtx *gin.Context) error {
	if bindingErr := ginCtx.ShouldBindBodyWith(&controller.apiInput, binding.JSON); bindingErr != nil {
		return bindingErr
	}
	controller.apiData.UserId = globalUtility.ConvertValueToInt(controller.apiInput.UserId)
	expirationTime, expirationTimeErr := time.Parse("2006-01-02", controller.apiInput.ExpirationTime)
	if expirationTimeErr != nil {
		return expirationTimeErr
	}
	controller.apiData.ExpirationTime = expirationTime
	validationStartTime, validationStartTimeErr := time.Parse("2006-01-02", controller.apiInput.ValidationStartTime)
	if validationStartTimeErr != nil {
		return validationStartTimeErr
	}
	controller.apiData.ValidationStartTime = validationStartTime
	return nil
}

func (controller *CreateTokenController) validateInput() error {
	if controller.apiData.UserId == 0 {
		return fmt.Errorf("invalid user_id")
	}
	if controller.apiData.ExpirationTime.Before(controller.apiData.ValidationStartTime) {
		return fmt.Errorf("expiration_time cant be before validation_start_time")
	}
	if controller.apiData.ExpirationTime.Sub(controller.apiData.ValidationStartTime) > 7*24*time.Hour {
		return fmt.Errorf("token can be valid for 7 days at most")
	}
	return nil
}
