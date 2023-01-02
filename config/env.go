package config

import (
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	BinanceEndpoint string `name:"BINANCE_ENDPOINT"`
	APIKey          string `name:"API_KEY"`
	APISecret       string `name:"API_SECRET"`

	Database *Database `name:"DATABASE"`
}

type Database struct {
	Host     string `name:"HOST"`
	Name     string `name:"DATABASE"`
	Username string `name:"USERNAME"`
	Password string `name:"PASSWORD"`
	Port     string `name:"PORT"`
}

var Environment Config

func LoadEnvironment() (err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}

	Environment.Database = &Database{}
	scanStruct(&Environment)

	return
}

func scanStruct(envVar any, name ...string) {
	typeOfENV := reflect.TypeOf(envVar)
	valueOfENV := reflect.ValueOf(envVar)

	for i := 0; i < typeOfENV.Elem().NumField(); i++ {
		envName := typeOfENV.Elem().Field(i).Tag.Get("name")

		if len(name) > 0 {
			envName = name[0] + "_" + envName
		}

		value := os.Getenv(envName)

		if typeOfENV.Elem().Field(i).Type.Kind() == reflect.Pointer {
			temp := valueOfENV.Elem().Field(i).Interface()
			scanStruct(temp, envName)
			continue
		}

		valueOfENV.Elem().Field(i).SetString(value)
	}
}
