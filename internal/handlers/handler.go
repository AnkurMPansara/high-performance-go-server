package handler

import (
	"backend-server/internal/modules/service/controllers/createTokenController"
	getGreetingsController "backend-server/internal/modules/service/controllers/getGreetings"

	"github.com/gin-gonic/gin"
)

func RouteRequests(ginEngine *gin.Engine) {
	serviceRequestHandler := ginEngine.Group("/service")
	serviceRequestHandler.POST("/getGreetings", getGreetingsController.GetGreetings)
	serviceRequestHandler.POST("/getGreetings/", getGreetingsController.GetGreetings)
	serviceRequestHandler.POST("/createToken", createTokenController.CreateToken)
	serviceRequestHandler.POST("/createToken/", createTokenController.CreateToken)
}