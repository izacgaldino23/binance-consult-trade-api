package candle

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/model"
	"github.com/izacgaldino23/binance-consult-trade-api/persist"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

var (
	CashBuyTotal                  float64 = 1000 // USD
	BuyPercent                            = 100
	SoldPercent                           = 100
	Bought                                = false
	sleepTo                       *time.Time
	boughtPrice                   float64
	CashSoldTotal                 float64 // BTC
	NumTransactions, sellAttempts int
)

const (
	maxTransactions = 1
	maxSellAttempts = 3
)

func BuyActive(price float64) (bool, error) {
	if Bought {
		return Bought, nil
	}

	// Calculate how much i will expand buying
	buy := CashBuyTotal * float64(BuyPercent) / 100 // USD

	// Calculate how much i will receive
	newCash := buy / price // BTC

	// Removing unnecessary decimals
	newCash = formatNumber(newCash, 8)

	// Remove from buy value this cents
	buy = newCash * price

	// Withdraw from our wallet
	CashBuyTotal -= buy

	// Added newCash bought
	CashSoldTotal += newCash

	boughtPrice = price

	// Logg("BUY ", price)
	tradeLogg("BUY ", "USD", "BTC", buy, newCash, price)

	Bought = true

	if err := SaveTrade(&model.Trade{
		Type:      model.TradeTypeBuy,
		FromName:  "USD",
		ToName:    "BTC",
		FromValue: buy,
		ToValue:   newCash,
		Price:     price,
		RSI:       RSI,
	}); err != nil {
		return Bought, err
	}

	NumTransactions++

	return Bought, nil
}

func SellActive(price float64, stopChan chan bool) (err error) {
	// validate if i have anything to sell
	if CashSoldTotal == 0 {
		return
	}

	// Validate if we can sell
	if !tryToSell(price) {
		return
	}

	sell := CashSoldTotal // BTC
	CashSoldTotal -= sell

	newCash := sell * price // USD
	CashBuyTotal += newCash

	tradeLogg("SELL", "BTC", "USD", sell, newCash, price)

	if NumTransactions == maxTransactions {
		Ticker.Stop()
		stopChan <- true
	}

	Bought = false
	sellAttempts = 0

	if err = SaveTrade(&model.Trade{
		Type:      model.TradeTypeSell,
		FromName:  "BTC",
		ToName:    "USD",
		FromValue: sell,
		ToValue:   newCash,
		Price:     price,
		RSI:       RSI,
	}); err != nil {
		return
	}

	return
}

func tryToSell(sellPrice float64) (sell bool) {
	// Verify if is sleeping period
	if sleepTo != nil && time.Now().Before(*sleepTo) {
		LoggInfo("Sleeping period")
		return
	}
	sleepTo = nil

	if sellAttempts == maxSellAttempts {
		LoggInfo("Sold for max attempts")
		return true
	}
	sellAttempts++

	if boughtPrice <= sellPrice {
		LoggInfo("Sold for good price")
		return true
	}

	LoggInfo("Not sold for lower value, attempt: ", sellAttempts)
	sleepTo = utils.GetTimePointer(time.Now().Add(time.Second * 10))

	return
}

func formatNumber(n float64, decimal int) float64 {
	exponencial := math.Pow(10, float64(decimal))
	temp, _ := math.Modf(n * exponencial)
	n = temp / exponencial

	return n
}

func LoggInfo(msg ...any) {
	temp := []any{model.BTCUSDT, " - "}
	temp = append(temp, msg...)
	temp = append(temp, " | ", time.Now().Format("02-01 15:04:05"))

	out := fmt.Sprint(temp...)
	outTransactions = append(outTransactions, out)

	fmt.Println(out)
}

func Logg(kind string, price float64) {
	out := fmt.Sprint(model.BTCUSDT, " - ", kind, " | USD: ", formatNumber(CashBuyTotal, 8), " | BTC: ", formatNumber(CashSoldTotal, 2), " | PRICE: ", price, " | RSI: ", RSI, " | ", time.Now().Format("02-01 15:04:05"))
	outTransactions = append(outTransactions, out)

	fmt.Println(out)
}

func tradeLogg(kind, expendType, boughtType string, expendValue, boughtValue, price float64) {
	out := fmt.Sprintf("%v - %v | %v: %2.8f | %v: %2.8f | PRICE: %.2f | RSI: %.2f | %v", model.BTCUSDT, kind, expendType, expendValue, boughtType, boughtValue, price, RSI, time.Now().Format("02-01 15:04:05"))

	outTransactions = append(outTransactions, out)

	fmt.Println(out)
}

func InitOrEndTransactions(price *float64, init bool) (err error) {
	text := model.TradeTypeInit

	if !init {
		text = model.TradeTypeEnd
	}

	Logg(strings.ToUpper(text), lastPrice)

	if err = SaveTrade(&model.Trade{
		Type:      text,
		FromName:  "USD",
		ToName:    "BTC",
		FromValue: CashBuyTotal,
		ToValue:   CashSoldTotal,
		Price:     *price,
		RSI:       RSI,
	}); err != nil {
		return
	}

	return
}

func SaveTrade(trade *model.Trade) (err error) {
	tx, err := config.NewTransaction(false)
	if err != nil {
		return
	}
	defer tx.Rollback()

	persistTrade := persist.TradePS{TX: tx}

	if _, err = persistTrade.AddTrade(trade); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
