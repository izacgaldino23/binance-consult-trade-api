package model

import (
	"reflect"
	"strconv"
	"time"
)

const (
	BTCUSDT = "BTCUSDT"
)

type Candle struct {
	ID         int64     `pos:"-" collumn:"id" ignoreInsertUpdate:"true"`
	OpenTime   time.Time `pos:"0" collumn:"open_time"`
	CloseTime  time.Time `pos:"6" collumn:"close_time"`
	OpenPrice  float64   `pos:"1" collumn:"open_price"`
	ClosePrice float64   `pos:"4" collumn:"close_price"`
	LowPrice   float64   `pos:"3" collumn:"low_price"`
	HighPrice  float64   `pos:"2" collumn:"high_price"`
	Volume     float64   `pos:"5" collumn:"volume"`
	Symbol     string    `pos:"-" collumn:"symbol"`
}

func (c *Candle) TableName() string {
	return "candle"
}

func (c *Candle) ArrayToStruct(values []interface{}) {
	typeOfENV := reflect.TypeOf(c)
	valueOfENV := reflect.ValueOf(c)

	for j := 0; j < typeOfENV.Elem().NumField(); j++ {
		posArray := typeOfENV.Elem().Field(j).Tag.Get("pos")

		if posArray == "-" {
			continue
		}

		posNumber, _ := strconv.ParseInt(posArray, 10, 64)

		kind := typeOfENV.Elem().Field(j).Type.Kind()

		switch kind {
		case reflect.Float32, reflect.Float64:
			value, _ := strconv.ParseFloat(values[posNumber].(string), 64)
			valueOfENV.Elem().Field(j).SetFloat(value)
		case reflect.Struct: // is time
			if typeOfENV.Elem().Field(j).Type == reflect.TypeOf(time.Now()) {
				value := values[posNumber].(float64) / 1000
				valueOfENV.Elem().Field(j).Set(reflect.ValueOf(time.Unix(int64(value), 0)))
			}
		}
	}
}
