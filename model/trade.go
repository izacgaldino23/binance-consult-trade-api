package model

import "time"

type Trade struct {
	ID        int64     `collumn:"id"`
	Type      string    `collumn:"type"`
	FromName  string    `collumn:"from_name"`
	ToName    string    `collumn:"to_name"`
	FromValue float64   `collumn:"from_value"`
	ToValue   float64   `collumn:"to_value"`
	Price     float64   `collumn:"price"`
	RSI       float64   `collumn:"rsi"`
	CreatedAt time.Time `collumn:"created_at"`
}

func (c *Trade) TableName() string {
	return "trade"
}
