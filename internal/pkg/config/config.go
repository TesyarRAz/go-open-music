package config

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type config struct {
	APP_DEBUG bool
	APP_PORT  int

	DB_DRIVER   string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USERNAME string
	DB_PASSWORD string

	JWT_ENCRYPT []byte
	JWT_SECRET  []byte
}

var AppConfig *config = &config{}

func Init() {
	AppConfig.APP_DEBUG, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	AppConfig.APP_PORT, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	AppConfig.DB_DRIVER = os.Getenv("DB_DRIVER")
	AppConfig.DB_HOST = os.Getenv("DB_HOST")
	AppConfig.DB_PORT = os.Getenv("DB_PORT")
	AppConfig.DB_NAME = os.Getenv("DB_NAME")
	AppConfig.DB_USERNAME = os.Getenv("DB_USERNAME")
	AppConfig.DB_PASSWORD = os.Getenv("DB_PASSWORD")

	AppConfig.JWT_ENCRYPT = []byte(os.Getenv("JWT_ENCRYPT"))
	AppConfig.JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	if !AppConfig.APP_DEBUG {
		gin.SetMode(gin.ReleaseMode)
	}
}
