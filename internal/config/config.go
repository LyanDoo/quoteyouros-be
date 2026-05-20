package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	CORS     CORSConfig
	App      AppConfig
}

type ServerConfig struct {
	Port int
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type EmailConfig struct {
	Provider      string
	ResendAPIKey  string
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	FromEmail     string
}

type CORSConfig struct {
	Origins []string
}

type AppConfig struct {
	AdminEmail string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	viper.SetDefault("SERVER_PORT", 8000)
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "quoteyouros")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("EMAIL_PROVIDER", "resend")
	viper.SetDefault("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")

	expDuration, err := time.ParseDuration(viper.GetString("JWT_EXPIRATION"))
	if err != nil {
		log.Fatalf("Invalid JWT expiration format: %v", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: viper.GetInt("SERVER_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: expDuration,
		},
		Email: EmailConfig{
			Provider:     viper.GetString("EMAIL_PROVIDER"),
			ResendAPIKey: viper.GetString("RESEND_API_KEY"),
			SMTPHost:     viper.GetString("SMTP_HOST"),
			SMTPPort:     viper.GetInt("SMTP_PORT"),
			SMTPUser:     viper.GetString("SMTP_USER"),
			SMTPPassword: viper.GetString("SMTP_PASSWORD"),
			FromEmail:    viper.GetString("EMAIL_FROM"),
		},
		CORS: CORSConfig{
			Origins: viper.GetStringSlice("CORS_ORIGINS"),
		},
		App: AppConfig{
			AdminEmail: viper.GetString("ADMIN_EMAIL"),
		},
	}
}
