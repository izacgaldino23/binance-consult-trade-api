package binance

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func SocketStart() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var subscribed bool

	// listenKey, err := binance.GetListenKey()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _ = listenKey

	u := url.URL{
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   "/stream",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println(string(message))
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if !subscribed {
				message := struct {
					ID     int      `json:"id"`
					Method string   `json:"method"`
					Params []string `json:"params"`
				}{
					1,
					"SUBSCRIBE",
					[]string{"btcusdt@kline_1m"},
				}
				err := conn.WriteJSON(message)
				if err != nil {
					log.Fatalln(err)
					return
				}
				subscribed = true
			} else {

			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
