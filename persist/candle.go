package persist

import (
	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/model"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

type CandlePS struct {
	TX *config.Transaction
}

func (c *CandlePS) AddCandle(candle *model.Candle) (id int64, err error) {
	valueMap, err := utils.ParseStructToMap(candle)
	if err != nil {
		return
	}

	if err = c.TX.Builder.
		Insert(candle.TableName()).
		SetMap(valueMap).
		Suffix("RETURNING id").
		Scan(&id); err != nil {
		return
	}

	return
}
