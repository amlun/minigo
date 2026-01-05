package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Init initializes Viper configuration with sensible defaults.
// It supports environment variables and optional config file (config.yaml/.json/.toml).
func Init() {
	// Load .env file first
	if err := gotenv.Load(); err != nil {
		log.Printf(".env file not found or failed to load: %v", err)
	}

	// Defaults
	viper.SetDefault("ENV", "prod")
	viper.SetDefault("PORT", "8808")
	viper.SetDefault("DB_DSN", "postgres://postgres:postgres@localhost:5432/chuangke?sslmode=disable")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("JWT_SECRET", "dev_secret_change_me")
	viper.SetDefault("JWT_EXPIRE_DURATION", "1h")

	// OSS配置
	viper.SetDefault("OSS_ENDPOINT", "")
	viper.SetDefault("OSS_ACCESS_KEY_ID", "")
	viper.SetDefault("OSS_ACCESS_KEY_SECRET", "")
	viper.SetDefault("OSS_BUCKET_NAME", "")
	viper.SetDefault("OSS_REGION", "")
	viper.SetDefault("OSS_TOKEN_EXPIRE_SECONDS", 3600)

	// Enable env overrides
	viper.AutomaticEnv()

	// Optional config file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./cmd/server")
	viper.AddConfigPath("./backend")
	viper.AddConfigPath("./backend/cmd/server")
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("loaded config file: %s", viper.ConfigFileUsed())
	}
}

// Helpers

func GetPort() string      { return viper.GetString("PORT") }
func GetDBDsn() string     { return viper.GetString("DB_DSN") }
func GetLogLevel() string  { return viper.GetString("LOG_LEVEL") }
func GetJWTSecret() string { return viper.GetString("JWT_SECRET") }
func GetJWTExpireDuration() time.Duration {
	d := viper.GetString("JWT_EXPIRE_DURATION")
	if dur, err := time.ParseDuration(d); err == nil {
		return dur
	}
	return time.Hour
}

func GetEnv() string  { return viper.GetString("ENV") }
func IsDevEnv() bool  { return GetEnv() == "dev" }
func IsProdEnv() bool { return GetEnv() == "prod" }

func GetOSSEndpoint() string        { return viper.GetString("OSS_ENDPOINT") }
func GetOSSAccessKeyID() string     { return viper.GetString("OSS_ACCESS_KEY_ID") }
func GetOSSAccessKeySecret() string { return viper.GetString("OSS_ACCESS_KEY_SECRET") }
func GetOSSBucketName() string      { return viper.GetString("OSS_BUCKET_NAME") }
func GetOSSRegion() string          { return viper.GetString("OSS_REGION") }
func GetOSSTokenExpireSeconds() int { return viper.GetInt("OSS_TOKEN_EXPIRE_SECONDS") }
