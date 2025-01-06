package utils

import (
	"go-kafka-producer/logger"
	"os"

	"github.com/joho/godotenv"
)

var DotEnvUtil IDotEnvUtil = NewDotEnvUtil()

type IDotEnvUtil interface {
	GetEnvVariable(key string) string
}

type dotEnvUtilStruct struct{}

func NewDotEnvUtil() IDotEnvUtil {
	return &dotEnvUtilStruct{}
}

func (d *dotEnvUtilStruct) GetEnvVariable(key string) string {
	err := godotenv.Load(`.env`)

	if err != nil {
		logger.Logger.Fatal("Error loading .env file in dotEnvUtil.go " + err.Error())
	}

	return os.Getenv(key)
}
