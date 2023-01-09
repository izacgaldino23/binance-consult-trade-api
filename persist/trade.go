package persist

import (
	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/model"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

type TradePS struct {
	TX *config.Transaction
}

func (c *TradePS) AddProcess(process *model.Process) (id int64, err error) {
	valueMap, err := utils.ParseStructToMap(process)
	if err != nil {
		return
	}

	if err = c.TX.Builder.
		Insert(process.TableName()).
		SetMap(valueMap).
		Suffix("RETURNING id").
		Scan(&id); err != nil {
		return
	}

	return
}

func (c *TradePS) AddTrade(trade *model.Trade) (id int64, err error) {
	valueMap, err := utils.ParseStructToMap(trade)
	if err != nil {
		return
	}

	if err = c.TX.Builder.
		Insert(trade.TableName()).
		SetMap(valueMap).
		Suffix("RETURNING id").
		Scan(&id); err != nil {
		return
	}

	return
}
