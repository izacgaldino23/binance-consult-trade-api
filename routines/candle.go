package routines

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/izacgaldino23/binance-consult-trade-api/binance"
	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/model"
	"github.com/izacgaldino23/binance-consult-trade-api/persist"
)

const (
	interval = time.Second * 5
)

var (
	errChan = make(chan error, 15)
)

func CandleWatch() {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			candleUpdate(errChan)
		case err := <-errChan:
			log.Fatal(err)
		}
	}
}

// 1. Get data from binance
// 2. Convert to our struct
// 3. Persist on database
// 4. Call trade funcs

func candleUpdate(errChan chan error) {
	var (
		result  string
		err     error
		candles []model.Candle
		symbol  = model.BTCUSDT
	)

	// Get candles from binance
	result, err = binance.GetCandle(symbol, 10)
	if err != nil {
		errChan <- err
		return
	}

	// convert candles to struct
	if candles, err = convertToStruct(result, symbol); err != nil {
		errChan <- err
		return
	}

	// Persist candles on database
	if err = saveCandles(candles); err != nil {
		errChan <- err
		return
	}
}

func convertToStruct(body, symbol string) (candles []model.Candle, err error) {
	candles = make([]model.Candle, 0)

	body = strings.Replace(body, "[[", "", -1)
	body = strings.Replace(body, "]]", "", -1)
	parts := strings.Split(body, "],[")

	// loop each candle string array
	for i := range parts {
		candle := model.Candle{}
		toArray := make([]interface{}, 0)
		value := fmt.Sprintf("[%v]", parts[i])

		if err = json.Unmarshal([]byte(value), &toArray); err != nil {
			return
		}

		candle.ArrayToStruct(toArray)
		candle.Symbol = symbol
		candles = append(candles, candle)
	}

	return
}

func saveCandles(candles []model.Candle) (err error) {
	tx, err := config.NewTransaction(false)
	defer tx.Rollback()

	persist := persist.CandlePS{TX: tx}

	for i := range candles {
		_, err = persist.AddCandle(&candles[i])

		if err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
