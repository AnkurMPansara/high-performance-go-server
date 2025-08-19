package getGreetingsController

import (
	"backend-server/internal/modules/service/models/getGreetingsModel"
	"backend-server/utilities/globalUtility"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GetGreetingsController struct {
	apiData getGreetingsModel.ApiData
	apiInput getGreetingsModel.ApiInput
	apiResponse getGreetingsModel.ApiResponse
}

func GetGreetings(ginCtx *gin.Context) {
	controller := &GetGreetingsController{}
	controller.apiData.StartTime = time.Now()
	controller.apiData.Code = http.StatusOK
	controller.perform(ginCtx)
}

func (controller *GetGreetingsController) perform(ginCtx *gin.Context) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			errorMsg := "Issue in getGreetings API: " + fmt.Sprintf("%v", panicErr)
			controller.apiData.Error = errorMsg
			controller.apiData.Code = http.StatusInternalServerError
			controller.returnApiResp(ginCtx)
		}
		getGreetingsModel.CreateLogs(ginCtx, controller.apiInput, controller.apiData)
	}()
	if bindingErr := controller.bindInput(ginCtx); bindingErr != nil {
		return
	}
	controller.apiData.ChatResponse = getGreetingsModel.FetchGreetings(controller.apiData.Reply)
	controller.apiResponse = getGreetingsModel.ApiResponse{
		Code: http.StatusOK,
		Status: "Success",
		Respose: controller.apiData.ChatResponse,
		Error: "",
	}
	controller.returnApiResp(ginCtx)
}

func (controller *GetGreetingsController) returnApiResp(ginCtx *gin.Context) {
	if !ginCtx.Writer.Written() {
		ginCtx.JSON(controller.apiResponse.Code, controller.apiResponse)
	}
}

func (controller *GetGreetingsController) bindInput(ginCtx *gin.Context) error {
	if bindingErr := ginCtx.ShouldBindBodyWith(&controller.apiInput, binding.JSON); bindingErr != nil {
		return bindingErr
	}
	controller.apiData.Reply = globalUtility.ConvertValueToString(controller.apiInput.Reply)
	return nil
}