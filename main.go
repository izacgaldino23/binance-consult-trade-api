package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/izacgaldino23/binance-consult-trade-api/binance"
	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/routines/candle"
	// "nhooyr.io/websocket"
)

func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

func main() {
	// LOAD ENVIRONMENT
	err := config.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// OPEN CONNECTIONS
	err = config.OpenConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Start app route
	app := fiber.New()

	app.Get("/ping", binance.Ping)
	// app.Get("/candle", binance.GetCandle)

	// go socket()
	// go binance.SocketStart()

	candle.CandleWatch()

	// Listes port 3000 for routes
	// err = app.Listen(":3000")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func socket() {
	// listenKey, err := binance.GetListenKey()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _ = listenKey

	// // conn, _, err := websocket.Dial(context.Background(), fmt.Sprintf("wss://stream.binance.com:9443/ws/%v", listenKey), nil)
	// conn, _, err := websocket.Dial(context.Background(), "wss://stream.binance.com:9443/stream", nil)
	// if err != nil {
	// 	// body, _ := ioutil.ReadAll(res.Body)
	// 	// log.Println(string(body))
	// 	log.Fatal(err)
	// }
	// defer conn.Close(http.StatusOK, "Connection finished")

	// message := struct {
	// 	ID     int      `json:"id"`
	// 	Method string   `json:"method"`
	// 	Params []string `json:"params"`
	// }{
	// 	1,
	// 	"LIST_SUBSCRIPTIONS",
	// 	[]string{"btcusdt@bookTicker", "bnbbtc@bookTicker"},
	// }

	// if err = wsjson.Write(context.Background(), conn, message); err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// for {
	// 	_, msg, err := conn.Read(context.Background())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		break
	// 	}

	// 	log.Println(string(msg))
	// }
}
