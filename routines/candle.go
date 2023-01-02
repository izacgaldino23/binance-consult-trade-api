package routines

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/izacgaldino23/binance-consult-trade-api/binance"
	"github.com/izacgaldino23/binance-consult-trade-api/model"
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
		result string
		err    error
	)

	result, err = binance.GetCandle(model.BTCUSDT, 10)
	if err != nil {
		errChan <- err
		return
	}

	if err = convertToStruct(result); err != nil {
		errChan <- err
		return
	}
}

func convertToStruct(body string) (err error) {
	var (
	// candles = make([]model.Candle, 0)
	// toArray = make(map[int64][]string)
	)

	body = strings.Replace(body, "[[", "", -1)
	body = strings.Replace(body, "]]", "", -1)
	parts := strings.Split(body, "],[")

	// loop each candle string array
	for i := range parts {
		// candle := model.Candle{}
		toArray := make([]interface{}, 0)
		value := fmt.Sprintf("[%v]", parts[i])

		if err = json.Unmarshal([]byte(value), &toArray); err != nil {
			return
		}

		fmt.Println(value)
	}

	// fmt.Println(toArray)

	return
}