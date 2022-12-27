package binance

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

func GetListenKey() (listenKey string, err error) {
	res, _, err := utils.Post(&utils.Request{
		URL: "/v3/userDataStream",
	})
	if err != nil {
		return
	}

	var (
		decoded = make(map[string]string, 0)
	)

	if err = json.Unmarshal(res, &decoded); err != nil {
		return
	}

	return
}

func Ping(c *fiber.Ctx) error {
	req := utils.Request{
		URL: "/v3/exchangeInfo",
	}

	body, _, err := utils.Get(req)
	if err != nil {
		log.Fatal(err)
	}

	return c.Send(body)
}

func GetCandle(c *fiber.Ctx) error {
	req := utils.Request{
		URL:   "/v3/klines",
		Query: &utils.QueryParamList{},
	}

	req.Query.
		AddParam("interval", "1m").
		AddParam("limit", "10").
		AddParam("symbol", "BTCUSDT")

	// req.ParseParams(c)

	body, _, err := utils.Get(req)
	if err != nil {
		log.Fatal(err)
	}

	return c.Send(body)
}
