package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	RabbitMQ   RabbitMQConfig
	Kafka      KafkaConfig
	Keycloak   KeycloakConfig
	Monitoring MonitoringConfig
	JWT        JWTConfig
	RateLimit  RateLimitConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	Version string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

type KafkaConfig struct {
	Brokers  []string
	GroupID  string
	ClientID string
}

type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
}

type MonitoringConfig struct {
	PrometheusNamespace string
	LogLevel            string
	LogFormat           string
}

type JWTConfig struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

// Load reads configuration from file and environment
func Load() (*Config, error) {
	// Set defaults first
	setDefaults()

	// Try to read from .env file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		// If .env doesn't exist, that's okay, we'll use environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// But if there's another error, log it
			fmt.Printf("Error reading .env file: %v\n", err)
		}
	}

	// Always read from environment variables (overrides .env file)
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	cfg := &Config{
		App: AppConfig{
			Name:    viper.GetString("APP_NAME"),
			Env:     viper.GetString("APP_ENV"),
			Port:    viper.GetString("APP_PORT"),
			Version: viper.GetString("APP_VERSION"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			User:     viper.GetString("RABBITMQ_USER"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
			VHost:    viper.GetString("RABBITMQ_VHOST"),
		},
		Kafka: KafkaConfig{
			Brokers:  viper.GetStringSlice("KAFKA_BROKERS"),
			GroupID:  viper.GetString("KAFKA_GROUP_ID"),
			ClientID: viper.GetString("KAFKA_CLIENT_ID"),
		},
		Keycloak: KeycloakConfig{
			URL:          viper.GetString("KEYCLOAK_URL"),
			Realm:        viper.GetString("KEYCLOAK_REALM"),
			ClientID:     viper.GetString("KEYCLOAK_CLIENT_ID"),
			ClientSecret: viper.GetString("KEYCLOAK_CLIENT_SECRET"),
		},
		Monitoring: MonitoringConfig{
			PrometheusNamespace: viper.GetString("PROMETHEUS_NAMESPACE"),
			LogLevel:            viper.GetString("LOG_LEVEL"),
			LogFormat:           viper.GetString("LOG_FORMAT"),
		},
		JWT: JWTConfig{
			AccessTokenExpiry:  viper.GetDuration("JWT_ACCESS_TOKEN_EXPIRY"),
			RefreshTokenExpiry: viper.GetDuration("JWT_REFRESH_TOKEN_EXPIRY"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: viper.GetInt("RATE_LIMIT_REQUESTS_PER_MINUTE"),
			Burst:             viper.GetInt("RATE_LIMIT_BURST"),
		},
	}

	// Special handling for Kafka brokers if it's a single string
	if len(cfg.Kafka.Brokers) == 0 {
		brokerString := viper.GetString("KAFKA_BROKERS")
		if brokerString != "" {
			cfg.Kafka.Brokers = []string{brokerString}
		}
	}

	return cfg, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("APP_NAME", "nucleus")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_VERSION", "1.0.0")

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)

	// RabbitMQ defaults
	viper.SetDefault("RABBITMQ_HOST", "localhost")
	viper.SetDefault("RABBITMQ_PORT", "5672")
	viper.SetDefault("RABBITMQ_VHOST", "/")

	// Monitoring defaults
	viper.SetDefault("PROMETHEUS_NAMESPACE", "nucleus")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	// JWT defaults
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRY", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRY", "7d")

	// Rate limit defaults
	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 60)
	viper.SetDefault("RATE_LIMIT_BURST", 10)
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required (current: '%s')", c.Database.Host)
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required (current: '%s')", c.Database.Name)
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}
	return nil
}
