package utils

import (
	"fmt"
	"math"

	"github.com/izacgaldino23/binance-consult-trade-api/model"
)

const (
	candleLimitCalculationRSI = 14
)

func CalculateIndicatorRSI(
	candles []model.Candle,
	lastPrice *float64,
	firstInteration *bool,
	preAvgGain *float64,
	preAvgLoss *float64,
	lastShowRSI *float64,
	RSI *float64,
) {
	var (
		gains              = make(Slice[float64], 0)
		losses             = make(Slice[float64], 0)
		avgGain, avgLoss   float64
		lastGain, lastLoss float64
	)

	// Get prices and gains and losses
	for i := range candles {
		price := &candles[i].ClosePrice

		// output = append(output, make(map[string]float64))

		if *lastPrice != 0 {
			difference := math.Ceil((*price-*lastPrice)*100) / 100

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

		*lastPrice = *price

		// output[len(output)-1]["price"] = *price
		// output[len(output)-1]["gain"] = lastGain
		// output[len(output)-1]["loss"] = lastLoss
		// output[len(output)-1]["avg_gain"] = 0
		// output[len(output)-1]["avg_loss"] = 0
	}

	if *firstInteration {
		// Calculate avg for gains and losses
		gains.Each(func(i int, v *float64) {
			avgGain += *v
		})

		losses.Each(func(i int, v *float64) {
			avgLoss += *v
		})

		avgGain = avgGain / candleLimitCalculationRSI
		avgLoss = avgLoss / candleLimitCalculationRSI

		*firstInteration = false
	} else {
		avgGain = (*preAvgGain*(candleLimitCalculationRSI-1) + lastGain) / candleLimitCalculationRSI
		avgLoss = (*preAvgLoss*(candleLimitCalculationRSI-1) + lastLoss) / candleLimitCalculationRSI
	}

	// output[len(output)-1]["avgGain"] = avgGain
	// output[len(output)-1]["avgLoss"] = avgLoss

	// Set to preview avg
	*preAvgGain = avgGain
	*preAvgLoss = avgLoss

	// Calculate RS
	RS := avgGain / avgLoss

	//Calculate RSI
	*RSI = 100 - (100 / (1 + RS))
	*RSI = math.Ceil(100**RSI) / 100
	if *lastShowRSI == 0 {
		*lastShowRSI = *RSI
	}

	// result, _ := json.Marshal(output)
	// fmt.Println(string(result))
	if (*lastShowRSI < *RSI && *lastShowRSI+5 <= *RSI) || (*lastShowRSI > *RSI && *lastShowRSI-5 >= *RSI) {
		fmt.Println(model.BTCUSDT, *RSI)
		*lastShowRSI = *RSI
	}
}
