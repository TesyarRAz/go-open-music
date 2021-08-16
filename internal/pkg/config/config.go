package config

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type config struct {
	APP_DEBUG bool
	APP_PORT  int
	APP_URL   string

	DB_DRIVER   string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USERNAME string
	DB_PASSWORD string

	JWT_ENCRYPT []byte
	JWT_SECRET  []byte

	RABBITMQ_HOST     string
	RABBITMQ_PORT     string
	RABBITMQ_USERNAME string
	RABBITMQ_PASSWORD string

	REDIS_HOST string
	REDIS_PORT string

	MAILTRAP_HOST     string
	MAILTRAP_PORT     int
	MAILTRAP_USERNAME string
	MAILTRAP_PASSWORD string
}

var AppConfig *config = &config{}

func Init() {
	AppConfig.APP_DEBUG, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	AppConfig.APP_PORT, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	AppConfig.APP_URL = os.Getenv("APP_URL")

	AppConfig.DB_DRIVER = os.Getenv("DB_DRIVER")
	AppConfig.DB_HOST = os.Getenv("DB_HOST")
	AppConfig.DB_PORT = os.Getenv("DB_PORT")
	AppConfig.DB_NAME = os.Getenv("DB_NAME")
	AppConfig.DB_USERNAME = os.Getenv("DB_USERNAME")
	AppConfig.DB_PASSWORD = os.Getenv("DB_PASSWORD")

	AppConfig.JWT_ENCRYPT = []byte(os.Getenv("JWT_ENCRYPT"))
	AppConfig.JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	AppConfig.RABBITMQ_HOST = os.Getenv("RABBITMQ_HOST")
	AppConfig.RABBITMQ_PORT = os.Getenv("RABBITMQ_PORT")
	AppConfig.RABBITMQ_USERNAME = os.Getenv("RABBITMQ_USERNAME")
	AppConfig.RABBITMQ_PASSWORD = os.Getenv("RABBITMQ_PASSWORD")

	AppConfig.REDIS_HOST = os.Getenv("REDIS_HOST")
	AppConfig.REDIS_PORT = os.Getenv("REDIS_PORT")

	AppConfig.MAILTRAP_HOST = os.Getenv("MAILTRAP_HOST")
	AppConfig.MAILTRAP_PORT, _ = strconv.Atoi(os.Getenv("MAILTRAP_PORT"))
	AppConfig.MAILTRAP_USERNAME = os.Getenv("MAILTRAP_USERNAME")
	AppConfig.MAILTRAP_PASSWORD = os.Getenv("MAILTRAP_PASSWORD")

	if !AppConfig.APP_DEBUG {
		gin.SetMode(gin.ReleaseMode)
	}
}
