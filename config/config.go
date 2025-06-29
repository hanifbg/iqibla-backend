package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppPort           int    `mapstructure:"app_port"`
	DbDriver          string `mapstructure:"db_driver"`
	DbHost            string `mapstructure:"db_host"`
	DbPort            int    `mapstructure:"db_port"`
	DbUser            string `mapstructure:"db_user"`
	DbPassword        string `mapstructure:"db_password"`
	DbName            string `mapstructure:"db_name"`
	DbSSLMode         string `mapstructure:"db_sslmode"`
	HttpTimeout       int    `mapstructure:"http_timeout"`
	MidtransServerKey string `mapstructure:"midtrans_server_key"`
	MidtransClientKey string `mapstructure:"midtrans_client_key"`
	IsProduction      bool   `mapstructure:"is_production"`
	BaseURL           string `mapstructure:"base_url"`
}

var (
	lock      = &sync.Mutex{}
	appConfig *AppConfig
)

func GetConfig() (*AppConfig, error) {
	if appConfig != nil {
		return appConfig, nil
	}

	lock.Lock()
	defer lock.Unlock()

	if appConfig != nil {
		return appConfig, nil
	}

	appConfig, err := initConfig()
	return appConfig, err
}

func initConfig() (*AppConfig, error) {
	var finalConfig AppConfig

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName("app.config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		// Load from environment variables if config file is not found
		finalConfig.DbDriver = getEnvOrDefault("DB_DRIVER", "postgres")
		finalConfig.DbHost = getEnvOrDefault("DB_HOST", "localhost")
		finalConfig.DbPort = getEnvIntOrDefault("DB_PORT", 5432)
		finalConfig.DbUser = getEnvOrDefault("DB_USER", "postgres")
		finalConfig.DbPassword = getEnvOrDefault("DB_PASSWORD", "")
		finalConfig.DbName = getEnvOrDefault("DB_NAME", "iqibla_ecommerce")
		finalConfig.DbSSLMode = getEnvOrDefault("DB_SSLMODE", "disable")
		finalConfig.AppPort = getEnvIntOrDefault("APP_PORT", 8080)
		finalConfig.MidtransServerKey = getEnvOrDefault("MIDTRANS_SERVER_KEY", "")
		finalConfig.MidtransClientKey = getEnvOrDefault("MIDTRANS_CLIENT_KEY", "")
		finalConfig.IsProduction = getEnvBoolOrDefault("IS_PRODUCTION", false)
		return &finalConfig, nil
	}

	fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	finalConfig.AppPort = viper.GetInt("server.port")
	finalConfig.DbDriver = viper.GetString("database.driver")
	finalConfig.DbHost = viper.GetString("database.host")
	finalConfig.DbPort = viper.GetInt("database.port")
	finalConfig.DbUser = viper.GetString("database.username")
	finalConfig.DbPassword = viper.GetString("database.password")
	finalConfig.DbName = viper.GetString("database.dbname")
	finalConfig.DbSSLMode = viper.GetString("database.sslmode")
	finalConfig.MidtransServerKey = viper.GetString("payment.midtrans_server_key")
	finalConfig.MidtransClientKey = viper.GetString("payment.midtrans_client_key")
	finalConfig.IsProduction = viper.GetBool("payment.is_production")

	return &finalConfig, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d"); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}