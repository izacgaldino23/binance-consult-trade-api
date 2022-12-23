package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/izacgaldino23/binance-consult-trade-api/binance"
	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

func main() {
	// LOAD ENVIRONMENT
	err := config.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Start app route
	app := fiber.New()

	app.Get("/ping", binance.Ping)

	go socket()

	// Listes port 3000 for routes
	err = app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func socket() {
	listenKey, err := binance.GetListenKey()
	if err != nil {
		log.Fatal(err)
	}
	_ = listenKey

	// conn, _, err := websocket.Dial(context.Background(), fmt.Sprintf("wss://stream.binance.com:9443/ws/%v", listenKey), nil)
	conn, _, err := websocket.Dial(context.Background(), "wss://stream.binance.com:9443/stream", nil)
	if err != nil {
		// body, _ := ioutil.ReadAll(res.Body)
		// log.Println(string(body))
		log.Fatal(err)
	}
	defer conn.Close(http.StatusOK, "Connection finished")

	message := struct {
		ID     int      `json:"id"`
		Method string   `json:"method"`
		Params []string `json:"params"`
	}{
		1,
		"LIST_SUBSCRIPTIONS",
		[]string{"btcusdt@bookTicker", "bnbbtc@bookTicker"},
	}

	if err = wsjson.Write(context.Background(), conn, message); err != nil {
		log.Fatal(err)
		return
	}

	for {
		_, msg, err := conn.Read(context.Background())
		if err != nil {
			log.Fatal(err)
			break
		}

		log.Println(string(msg))
	}

	// cstDialer := websocket.Dialer{
	// 	Subprotocols:     []string{"p1", "p2"},
	// 	ReadBufferSize:   1024,
	// 	WriteBufferSize:  1024,
	// 	HandshakeTimeout: 30 * time.Second,
	// }

	// listenKey, err := binance.GetListenKey()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// conn, response, err := cstDialer.Dial(fmt.Sprintf("wss://stream.binance.com:443/ws/%v", listenKey), http.Header{
	// 	"X-MBX-APIKEY": []string{config.Environment.APIKey},
	// })
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// defer conn.Close()

	// log.Println(response, conn)
}
