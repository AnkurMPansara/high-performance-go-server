package globalUtility

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ConvertValueToInt(v interface{}) int {
	convertedValue := 0
	switch val := v.(type) {
	case int64:
		convertedValue = int(val)
	case int32: //rune and int32 are same
		convertedValue = int(val)
	case int16:
		convertedValue = int(val)
	case int8:
		convertedValue = int(val)
	case int:
		convertedValue = val
	case uint:
		convertedValue = int(val)
	case uint64:
		convertedValue = int(val)
	case uint32:
		convertedValue = int(val)
	case uint16:
		convertedValue = int(val)
	case uint8:
		convertedValue = int(val)
	case []byte:
		convertedValue, _ = strconv.Atoi(string(val))
	case string:
		convertedValue, _ = strconv.Atoi(string(val))
	case float32:
		convertedValue = int(val)
	case float64:
		convertedValue = int(val) 
	case json.Number:
		if num, err := val.Int64(); err == nil {
			convertedValue = int(num)
		} else if num, err := val.Float64(); err == nil {
			convertedValue = int(num)
		}
	case bool:
		if val {
			convertedValue = 1
		}
	case time.Duration:
		convertedValue = int(val)
	default:
		convertedValue = 0
	}
	return convertedValue
}

func ConvertValueToString(v interface{}) string {
	if v == nil {
		return ""
	}
	convertedValue := ""
	switch val := v.(type) {
	case map[string]interface{}:
		if jsonByteData, err := json.Marshal(val); err == nil {
			convertedValue = string(jsonByteData)
		}
	case []interface{}:
		if jsonByteData, err := json.Marshal(val); err == nil {
			convertedValue = string(jsonByteData)
		}
	case time.Time:
		convertedValue = val.Format("2006-01-02 15:04:05.000000000")
	case http.Header:
		if jsonByteData, err := json.Marshal(val); err == nil {
			convertedValue = string(jsonByteData)
		}
	case error:
		convertedValue = val.Error()
	default:
		convertedValue = fmt.Sprintf("%v", val)
	}
	return convertedValue
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