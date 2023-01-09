package candle

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/izacgaldino23/binance-consult-trade-api/binance"
	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/model"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

const (
	interval    = time.Second * 1
	candleLimit = 15
)

var (
	errChan                        = make(chan error, 15)
	stop                           = make(chan bool, 1)
	firstInteration                = true
	symbol                         = model.BTCUSDT
	outTransactions                = make([]string, 0)
	Ticker                         *time.Ticker
	RSI, lastShowRSI               float64
	lastCandleStartTime, processID int64
	preAvgGain, preAvgLoss         float64
	lastPrice                      float64
)

func CandleWatch() {
	Ticker = time.NewTicker(interval)

	for {
		select {
		case <-stop:
			errChan <- EndTransactions(&processID)

			fmt.Println("----------------------------")
			fmt.Println(strings.Join(outTransactions, "\n"))
			fmt.Println("----------------------------")

			return
		case <-Ticker.C:
			candleUpdate(errChan, stop)
		case err := <-errChan:
			log.Fatal(err)
		}
	}
}

// 1. Get data from binance
// 2. Convert to our struct
// 3. Persist on database
// 4. Calculate the indicators
// 5. Call trade funcs

func candleUpdate(errChan chan error, stopChan chan bool) {
	var (
		result  string
		err     error
		candles []model.Candle
	)

	// Get candles from binance
	result, err = binance.GetCandle(symbol, candleLimit, lastCandleStartTime)
	if err != nil {
		errChan <- err
		return
	}

	if result == "[]" {
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

	// Calculate the RSI indicators
	if err = calculateIndicators(candles, stopChan); err != nil {
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

	// persist := persist.CandlePS{TX: tx}

	for i := range candles {
		// _, err = persist.AddCandle(&candles[i])

		if err != nil {
			return
		}

		if candles[i].OpenTime.Unix()*1000 > int64(lastCandleStartTime) {
			lastCandleStartTime = (candles[i].OpenTime.Unix() * 1000) + 1
		}
	}

	// fmt.Printf("Saved %v candles\n", len(candles))

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func calculateIndicators(candles []model.Candle, stopChan chan bool) (err error) {
	utils.CalculateIndicatorRSI(candles, &lastPrice, &firstInteration, &preAvgGain, &preAvgLoss, &lastShowRSI, &RSI)

	if len(outTransactions) == 0 {
		if processID, err = InitTransactions(&lastPrice, &symbol); err != nil {
			return
		}
	}

	if RSI >= 70 {
		err = SellActive(lastPrice, &processID, stopChan)
	} else if RSI <= 30 {
		_, err = BuyActive(lastPrice, &processID)
	}

	return
}
