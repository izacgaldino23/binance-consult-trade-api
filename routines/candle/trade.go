package candle

import (
	"fmt"
	"math"
	"time"

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

	// Removing unnecessary decimals
	newCash = formatNumber(newCash, 2)

	// Remove from buy value this cents
	buy = newCash / price

	// Withdraw from our wallet
	CashBuyTotal -= buy

	// Added newCash bought
	CashSoldTotal += newCash

	// Logg("BUY ", price)
	tradeLogg("BUY ", "BTC", "USD", buy, newCash, price)

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

	// Logg("SELL", price)
	tradeLogg("SELL", "USD", "BTC", sell, newCash, price)

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
	out := fmt.Sprint(model.BTCUSDT, " - ", kind, " | BTC: ", formatNumber(CashBuyTotal, 8), " | DOLLAR: ", formatNumber(CashSoldTotal, 2), " | PRICE: ", price, " | RSI: ", RSI, " | ", time.Now().Format("02-01 15:04:05"))
	outTransactions = append(outTransactions, out)

	fmt.Println(out)
}

func tradeLogg(kind, expendType, boughtType string, expendValue, boughtValue, price float64) {
	out := fmt.Sprintf("%v - %v | %v: %2.8f | %v: %2.8f | PRICE: %.2f | RSI: %.2f | %v", model.BTCUSDT, kind, expendType, expendValue, boughtType, boughtValue, price, RSI, time.Now().Format("02-01 15:04:05"))

	outTransactions = append(outTransactions, out)

	fmt.Println(out)
}
