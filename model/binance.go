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
	OpenTime   time.Time `pos:"0"`
	CloseTime  time.Time `pos:"6"`
	OpenPrice  float64   `pos:"1"`
	ClosePrice float64   `pos:"4"`
	LowPrice   float64   `pos:"3"`
	HighPrice  float64   `pos:"2"`
	Volume     float64   `pos:"5"`
}

func (c *Candle) ArrayToStruct(values []interface{}) {
	typeOfENV := reflect.TypeOf(c)
	valueOfENV := reflect.ValueOf(c)

	for j := 0; j < typeOfENV.Elem().NumField(); j++ {
		posArray := typeOfENV.Elem().Field(j).Tag.Get("pos")
		posNumber, _ := strconv.ParseInt(posArray, 10, 64)

		kind := typeOfENV.Elem().Field(j).Type.Kind()

		switch kind {
		case reflect.Float32, reflect.Float64:
			// valueOfENV.Elem().Field(j).SetFloat(values[posNumber].(float64))
			value, _ := strconv.ParseFloat(values[posNumber].(string), 64)
			valueOfENV.Elem().Field(j).SetFloat(value)
		case reflect.Struct: // is time
			if typeOfENV.Elem().Field(j).Type == reflect.TypeOf(time.Now()) {
				value := values[posNumber].(float64) / 1000
				valueOfENV.Elem().Field(j).Set(reflect.ValueOf(time.Unix(int64(value), 0)))
			}
			// case reflect.String:
			// 	value, _ := strconv.ParseFloat(values[posNumber].(string), 64)
			// 	valueOfENV.Elem().Field(j).SetFloat(value)
		}

		// valueOfENV.Elem().Field(j).SetString(values[posNumber])
	}
}
