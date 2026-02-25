package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv                 string        `mapstructure:"APP_ENV"`
	Port                   string        `mapstructure:"PORT"`
	DBHost                 string        `mapstructure:"DB_HOST"`
	DBPort                 string        `mapstructure:"DB_PORT"`
	DBUser                 string        `mapstructure:"DB_USER"`
	DBPassword             string        `mapstructure:"DB_PASSWORD"`
	DBName                 string        `mapstructure:"DB_NAME"`
	DBSSLMode              string        `mapstructure:"DB_SSLMODE"`
	JWTSecretKey           string        `mapstructure:"JWT_SECRET_KEY"`
	JWTAccessTokenExpHours time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXP_HOURS"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	var cfg Config
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: error reading config file: %v", err)
	}

	err = viper.Unmarshal(&cfg)

	// Set defaults if empty
	if cfg.Port == "" {
		cfg.Port = "8081"
	}
	if cfg.JWTSecretKey == "" {
		cfg.JWTSecretKey = "ops_super_secret_key"
	}
	if cfg.JWTAccessTokenExpHours == 0 {
		cfg.JWTAccessTokenExpHours = 24
	}

	return cfg, err
}
