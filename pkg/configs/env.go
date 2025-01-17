package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
	FromEmail              string
	FromEmailPassword      string
	FromEmailSMTP          string
	SMTPAddress            string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "toor"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTSecret:              getEnv("JWT_SECRET", "Uh3BnyZivL99alxVwRQpbjdkPFu2l9MnCSfgWn8HeXRPSlkXano7sdYYOwKhvpB+eq3mo9SRKpDTMdNqHOuQWA=="),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		FromEmail:              getEnv("FROM_EMAIL", ""),
		FromEmailPassword:      getEnv("FROM_EMAIL_PASSWORD", ""),
		FromEmailSMTP:          getEnv("FROM_EMAIL_SMTP", "smtp.gmail.com"),
		SMTPAddress:            getEnv("SMTP_ADDR", "smtp.gmail.com:587"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
