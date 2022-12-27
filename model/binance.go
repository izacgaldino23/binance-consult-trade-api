package model

import "time"

type Candle struct {
	OpenTime   *time.Time `pos:"0"`
	CloseTime  *time.Time `pos:"6"`
	OpenPrice  float64    `pos:"1"`
	ClosePrice float64    `pos:"4"`
	LowPrice   float64    `pos:"3"`
	HighPrice  float64    `pos:"2"`
	Volume     float64    `pos:"5"`
}
