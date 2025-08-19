package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

var configVar *viper.Viper

func LoadConfig() error {
	kolkataTimeZone, timeLoadErr := time.LoadLocation("Asia/Kolkata")
	if timeLoadErr != nil {
		return timeLoadErr
	}
	time.Local = kolkataTimeZone
	configVar = viper.New()
	configVar.AutomaticEnv()
	configVar.SetDefault("ENV", "DEV")
	if err := setConfigPath(configVar); err != nil {
		return err
	}
	return nil
}

func setConfigPath(config *viper.Viper) error {
	env := config.GetString("ENV")
	if env == "" {
		env = "DEV"
	}
	var configPath string
	var configName string
	switch env {
	case "DEV":
		configName = "config_dev"
	case "STG":
		configName = "config_stg"
	case "PROD":
		configName = "config_prod"
	default:
		configName = "config_dev"
	}
	if currentDirectory, err := os.Getwd(); err == nil {
		configPath = currentDirectory
	}
	configPath = filepath.Join(configPath, "config")

	config.AddConfigPath(configPath)
	config.SetConfigName("config_global")
	config.SetConfigType("yaml")
	if loadConfigError := config.ReadInConfig(); loadConfigError != nil {
		return fmt.Errorf("error loading global config: %w", loadConfigError)
	}
	config.SetConfigName(configName)
	config.SetConfigType("yaml")
	if loadConfigError := config.MergeInConfig(); loadConfigError != nil {
		return fmt.Errorf("error loading env config (%s): %w", configName, loadConfigError)
	}
	return nil
}

func GetConfigStringValue(key string) string {
	return configVar.GetString(key)
}

func GetConfigIntValue(key string) int {
	return configVar.GetInt(key)
}

func GetConfigMapValue(key string) map[string]interface{} {
	return configVar.GetStringMap(key)
}