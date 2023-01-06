package candle

import (
	"fmt"
	"math"

	"github.com/izacgaldino23/binance-consult-trade-api/model"
)

var (
	CashBuyTotal    float64 = 0.002 // BTC
	CashSoldTotal   float64 = 0     // Dollar
	BuyPercent              = 50
	SoldPercent             = 100
	NumTransactions         = 0
	Bought                  = false
)

const (
	maxTransactions = 1
)

func BuyActive(price float64) bool {
	if Bought {
		return Bought
	}

	// Calculate how much i will expand buying
	buy := CashBuyTotal * float64(BuyPercent) / 100 // BTC

	// Calculate how much i will receive
	newCash := buy * price // Dollar

	// Difference in cents
	difference := newCash - formatNumber(newCash, 2)

	// Remove from buy value this cents
	buy -= difference / price

	// Removing difference
	newCash -= difference

	// Withdraw from our wallet
	CashBuyTotal -= buy

	// Added newCash bought
	CashSoldTotal += newCash

	Logg("BUY ", price)

	NumTransactions++

	Bought = true

	return Bought
}

func SellActive(price float64, stopChan chan bool) (err error) {
	// validate if i have anything to sell
	if CashSoldTotal == 0 {
		return
	}

	// sell := CashSoldTotal * float64(SoldPercent) / 100 // Dollar
	sell := CashSoldTotal // Dollar
	CashSoldTotal -= sell

	newCash := sell / price // BTC
	CashBuyTotal += newCash

	Logg("SELL", price)

	if NumTransactions == maxTransactions {
		Ticker.Stop()
		stopChan <- true
		return
	}

	Bought = false

	return
}

func formatNumber(n float64, decimal int) float64 {
	exponencial := math.Pow(10, float64(decimal))
	temp, _ := math.Modf(n * exponencial)
	n = temp / exponencial

	return n
}

func Logg(kind string, price float64) {
	out := fmt.Sprint(model.BTCUSDT, " - ", kind, " | BTC: ", formatNumber(CashBuyTotal, 8), " | DOLLAR: ", formatNumber(CashSoldTotal, 2), " | PRICE: ", price, " | RSI: ", RSI, " | ")
	outTransactions = append(outTransactions, out)

	fmt.Println(out)
}
