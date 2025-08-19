package globalUtility

import (
	"backend-server/utilities/configuration"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var fileLock sync.Mutex

func CreateApplicationLogs(log map[string]interface{}) {
	finalLogData := make(map[string]interface{})
	for key, value := range log {
		switch val := value.(type) {
		case int:
			finalLogData[key] = val
		default:
			finalLogData[key] = ConvertValueToString(val)
		}
	}
	finalLogString := ConvertValueToString(finalLogData)
	fmt.Println(finalLogString)
	logFilePath := configuration.GetConfigStringValue("application_log_path")
	currentDirectory, _ := os.Getwd()
	logFilePath = filepath.Join(currentDirectory, logFilePath)
	logFileName := configuration.GetConfigStringValue("application_log_file") + "_" + time.Now().Format("2006-01-02") + ".log"
	if err := WriteInFile(finalLogString, filepath.Join(logFilePath, logFileName)); err != nil {
		fmt.Println("error loading writing logs: %w", err)
	} 
}

func CreateAccessLogs(log map[string]interface{}) {
	finalLogData := make(map[string]interface{})
	for key, value := range log {
		switch val := value.(type) {
		case int:
			finalLogData[key] = val
		default:
			finalLogData[key] = ConvertValueToString(val)
		}
	}
	finalLogString := ConvertValueToString(finalLogData)
	logFilePath := configuration.GetConfigStringValue("access_log_path")
	currentDirectory, _ := os.Getwd()
	logFilePath = filepath.Join(currentDirectory, logFilePath)
	logFileName := configuration.GetConfigStringValue("access_log_file") + "_" + time.Now().Format("2006-01-02") + ".log"
	if err := WriteInFile(finalLogString, filepath.Join(logFilePath, logFileName)); err != nil {
		fmt.Println("error loading writing logs: %w", err)
	} 
}

func WriteInFile(content string, filePath string) error {
	fileLock.Lock()
	file, fileOpenErr := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer func() {
		file.Close()
		fileLock.Unlock()
	}()
	_, writeErr := file.WriteString(fmt.Sprintf("%s\n", content))
	if writeErr != nil {
		return writeErr
	}
	return nil
}