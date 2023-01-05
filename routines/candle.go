package routines

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	errChan                = make(chan error, 15)
	firstInteration        = true
	output                 = make([]map[string]float64, 0)
	RSI                    float64
	lastCandleStartTime    int64
	preAvgGain, preAvgLoss float64
	lastPrice              float64
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
	result, err = binance.GetCandle(symbol, candleLimit, lastCandleStartTime)
	if err != nil {
		errChan <- err
		return
	}

	if result == "[]" {
		// fmt.Println("No new candles to add")
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
	calculateIndicatorRSI(candles)
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

func calculateIndicatorRSI(candles []model.Candle) {
	var (
		gains              = make(utils.Slice[float64], 0)
		losses             = make(utils.Slice[float64], 0)
		avgGain, avgLoss   float64
		lastGain, lastLoss float64
	)

	if firstInteration {
		// output[0] = append(output[0], "price")
		// output[1] = append(output[1], "gain")
		// output[2] = append(output[2], "loss")
		// output[3] = append(output[3], "avg_gain")
		// output[4] = append(output[4], "avg_loss")
	}

	// Get prices and gains and losses
	for i := range candles {
		price := &candles[i].ClosePrice

		output = append(output, make(map[string]float64))

		if lastPrice != 0 {
			difference := math.Ceil((*price-lastPrice)*100) / 100

			if difference > 0 {
				lastGain = difference
				gains = append(gains, lastGain)

				lastLoss = 0
				losses = append(losses, 0)
			} else {
				lastLoss = math.Abs(difference)
				losses = append(losses, lastLoss)

				lastGain = 0
				gains = append(gains, 0)
			}
		}

		lastPrice = *price

		output[len(output)-1]["price"] = *price
		output[len(output)-1]["gain"] = lastGain
		output[len(output)-1]["loss"] = lastLoss
		output[len(output)-1]["avg_gain"] = 0
		output[len(output)-1]["avg_loss"] = 0
	}

	if firstInteration {
		// Calculate avg for gains and losses
		gains.Each(func(i int, v *float64) {
			avgGain += *v
		})

		losses.Each(func(i int, v *float64) {
			avgLoss += *v
		})

		avgGain = avgGain / candleLimit
		avgLoss = avgLoss / candleLimit

		firstInteration = false
	} else {
		avgGain = (preAvgGain*(candleLimit-1) + lastGain) / candleLimit
		avgLoss = (preAvgLoss*(candleLimit-1) + lastLoss) / candleLimit
	}

	output[len(output)-1]["avgGain"] = avgGain
	output[len(output)-1]["avgLoss"] = avgLoss

	// Set to preview avg
	preAvgGain = avgGain
	preAvgLoss = avgLoss

	// Calculate RS
	RS := avgGain / avgLoss

	//Calculate RSI
	RSI = 100 - (100 / (1 + RS))

	// result, _ := json.Marshal(output)
	// fmt.Println(string(result))
	fmt.Println(model.BTCUSDT, RSI)
}
