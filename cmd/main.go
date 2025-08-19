package main

import (
	handler "backend-server/internal/handlers"
	"backend-server/internal/middlewares/accessLogger"
	"backend-server/internal/middlewares/authentication"
	"backend-server/internal/middlewares/customRecovery"
	"backend-server/utilities/configuration"
	"backend-server/utilities/globalUtility"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Server...")
	//start the server
	initiateServer()
}

func initiateServer() {
	defer func() {
		if panicRecovery := recover(); panicRecovery != nil {
			stackTrace := string(debug.Stack())
			errorLogData := make(map[string]interface{})
			errorLogData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_server_error")
			errorLogData["CODE"] = http.StatusInternalServerError
			errorLogData["ERROR"] = globalUtility.ConvertValueToString(panicRecovery)
			errorLogData["STACK"] = stackTrace
			globalUtility.CreateApplicationLogs(errorLogData)
		}
	}()

	configuration.LoadConfig()

	ginEngine := gin.New()
	ginEngine.Use(customRecovery.HandlePanic)
	ginEngine.Use(accessLogger.AccessLog)
	ginEngine.Use(authentication.AuthenticateRequest)

	handler.RouteRequests(ginEngine)

	serverPort := ""
	if len(os.Args) > 1 {
		serverPort = os.Args[1]
		if globalUtility.ConvertValueToInt(serverPort) == 0 {
			errorLogData := make(map[string]interface{})
			errorLogData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_server_error")
			errorLogData["CODE"] = http.StatusInternalServerError
			errorLogData["ERROR"] = "Invalid Port provided"
			globalUtility.CreateApplicationLogs(errorLogData)
			return
		}
	} else {
		errorLogData := make(map[string]interface{})
		errorLogData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_server_error")
		errorLogData["CODE"] = http.StatusInternalServerError
		errorLogData["ERROR"] = "No Port provided"
		globalUtility.CreateApplicationLogs(errorLogData)
		return
	}

	httpServer := &http.Server{
		Addr: ":" + serverPort,
		Handler: ginEngine,
	}

	startListening(httpServer)
}

func startListening(httpServer *http.Server) {
	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println("Some error occured the server ...")
		if err != http.ErrServerClosed {
			errorLogData := make(map[string]interface{})
			errorLogData["LOG_TYPE"] = configuration.GetConfigStringValue("log_type_server_error")
			errorLogData["CODE"] = http.StatusInternalServerError
			errorLogData["ERROR"] = fmt.Sprintf("Some error occured while starting the server: %w", err)
			globalUtility.CreateApplicationLogs(errorLogData)
		} else {
			fmt.Println("Server Closed")
		}
	}
}