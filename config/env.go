package config

import (
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	BinanceEndpoint string `name:"BINANCE_ENDPOINT"`
}

var Environment Config

func LoadEnvironment() (err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}

	temp := reflect.TypeOf(Environment)
	temp2 := reflect.ValueOf(&Environment)

	for i := 0; i < temp.NumField(); i++ {
		envName := temp.Field(i).Tag.Get("name")
		value := os.Getenv(envName)

		temp2.Elem().Field(i).SetString(value)
	}

	return
}
