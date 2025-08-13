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
	RajaOngkirAPIKey  string `mapstructure:"rajaongkir_api_key"`
	RajaOngkirBaseURL string `mapstructure:"rajaongkir_base_url"`

	// RajaOngkir caching configuration
	RajaOngkirCacheEnabled      bool   `mapstructure:"rajaongkir_cache_enabled"`
	RajaOngkirCacheTTLHours     int    `mapstructure:"rajaongkir_cache_ttl_hours"`
	RajaOngkirWarmupOnStartup   bool   `mapstructure:"rajaongkir_warmup_on_startup"`
	RajaOngkirWarmupTimeoutSecs int    `mapstructure:"rajaongkir_warmup_timeout_secs"`
	SMTPHost                    string `mapstructure:"smtp_host"`
	SMTPPort                    int    `mapstructure:"smtp_port"`
	SMTPUsername                string `mapstructure:"smtp_username"`
	SMTPPassword                string `mapstructure:"smtp_password"`
	SMTPFrom                    string `mapstructure:"smtp_from"`
	WhatsappConfig              WhatsappConfig
}

type WhatsappConfig struct {
	Host     string
	Username string
	Password string
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
		finalConfig.RajaOngkirAPIKey = getEnvOrDefault("RAJAONGKIR_API_KEY", "")
		finalConfig.RajaOngkirBaseURL = getEnvOrDefault("RAJAONGKIR_BASE_URL", "")
		finalConfig.SMTPHost = getEnvOrDefault("SMTP_HOST", "")
		finalConfig.SMTPPort = getEnvIntOrDefault("SMTP_PORT", 0)
		finalConfig.SMTPUsername = getEnvOrDefault("SMTP_USERNAME", "")
		finalConfig.SMTPPassword = getEnvOrDefault("SMTP_PASSWORD", "")
		finalConfig.SMTPFrom = getEnvOrDefault("SMTP_FROM", "")
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
	finalConfig.BaseURL = viper.GetString("base_url")
	finalConfig.HttpTimeout = viper.GetInt("http_timeout")
	finalConfig.RajaOngkirAPIKey = viper.GetString("shipping.rajaongkir_api_key")
	finalConfig.RajaOngkirBaseURL = viper.GetString("shipping.rajaongkir_base_url")

	// Load cache configuration
	finalConfig.RajaOngkirCacheEnabled = viper.GetBool("shipping.rajaongkir_cache_enabled")
	finalConfig.RajaOngkirCacheTTLHours = viper.GetInt("shipping.rajaongkir_cache_ttl_hours")
	finalConfig.RajaOngkirWarmupOnStartup = viper.GetBool("shipping.rajaongkir_warmup_on_startup")
	finalConfig.RajaOngkirWarmupTimeoutSecs = viper.GetInt("shipping.rajaongkir_warmup_timeout_secs")

	//email
	finalConfig.SMTPHost = viper.GetString("mail.host")
	finalConfig.SMTPPort = viper.GetInt("mail.port")
	finalConfig.SMTPUsername = viper.GetString("mail.username")
	finalConfig.SMTPPassword = viper.GetString("mail.password")
	finalConfig.SMTPFrom = viper.GetString("mail.sender_email")

	//WA
	finalConfig.WhatsappConfig.Host = viper.GetString("whatsapp.url")
	finalConfig.WhatsappConfig.Username = viper.GetString("whatsapp.username")
	finalConfig.WhatsappConfig.Password = viper.GetString("whatsapp.password")

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
