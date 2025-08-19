package sendMessageController

import (
	"backend-server/internal/modules/chat/models/sendMessageModel"
	"backend-server/utilities/globalUtility"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SendMessageController struct {
	apiData sendMessageModel.ApiData
	apiInput sendMessageModel.ApiInput
	apiResponse sendMessageModel.ApiResponse
}

func SendMessage(ginCtx *gin.Context) {
	controller := &SendMessageController{}
	controller.apiData.StartTime = time.Now()
	controller.apiData.Code = http.StatusOK
	controller.perform(ginCtx)
}

func (controller *SendMessageController) perform(ginCtx *gin.Context) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			errorMsg := "Issue in createToken API: " + fmt.Sprintf("%v", panicErr)
			controller.apiData.Error = errorMsg
			controller.apiData.Code = http.StatusInternalServerError
			controller.returnApiResp(ginCtx)
		}
		sendMessageModel.CreateLogs(ginCtx, controller.apiInput, controller.apiData)
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
}

func (controller *SendMessageController) bindInput(ginCtx *gin.Context) error {
	if bindingErr := ginCtx.ShouldBindBodyWith(&controller.apiInput, binding.JSON); bindingErr != nil {
		return bindingErr
	}
	controller.apiData.UserId = globalUtility.ConvertValueToInt(controller.apiInput.UserId)
	controller.apiData.Message = globalUtility.ConvertValueToString(controller.apiInput.Message)
	return nil
}

func (controller *SendMessageController) validateInput() error {
	if controller.apiData.UserId == 0 {
		return fmt.Errorf("invalid user_id")
	}
	if controller.apiData.Message == "" {
		return fmt.Errorf("message can not be empty")
	}
	return nil
}

func (controller *SendMessageController) returnApiResp(ginCtx *gin.Context) {
	if !ginCtx.Writer.Written() {
		ginCtx.JSON(controller.apiResponse.Code, controller.apiResponse)
	}
}