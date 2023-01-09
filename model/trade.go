package model

import "time"

const (
	TradeTypeBuy  = "buy"
	TradeTypeSell = "sell"
	TradeTypeInit = "init"
	TradeTypeEnd  = "end"
)

type Process struct {
	ID         int64     `collumn:"id" ignoreInsertUpdate:"true"`
	Symbol     string    `collumn:"symbol"`
	FromName   string    `collumn:"from_name"`
	ToName     string    `collumn:"to_name"`
	StartValue float64   `collumn:"start_value"`
	EndValue   float64   `collumn:"end_value"`
	StartPrice float64   `collumn:"start_price"`
	CreatedAt  time.Time `collumn:"created_at" ignoreInsertUpdate:"true"`
	UpdatedAt  time.Time `collumn:"created_at" ignoreInsertUpdate:"true"`
}

type Trade struct {
	ID        int64     `collumn:"id" ignoreInsertUpdate:"true"`
	ProcessID int64     `collumn:"process_id"`
	Type      string    `collumn:"type"`
	FromName  string    `collumn:"from_name"`
	ToName    string    `collumn:"to_name"`
	FromValue float64   `collumn:"from_value"`
	ToValue   float64   `collumn:"to_value"`
	Price     float64   `collumn:"price"`
	RSI       float64   `collumn:"rsi"`
	CreatedAt time.Time `collumn:"created_at" ignoreInsertUpdate:"true"`
}

func (c *Trade) TableName() string {
	return "trade"
}

func (c *Process) TableName() string {
	return "process"
}
