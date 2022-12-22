package main

import (
	"fmt"
	"log"

	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

const ()

func main() {
	// LOAD ENVIRONMENT
	err := config.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	req := utils.Request{
		URL: "/v3/ping",
	}

	_, status, err := utils.Get(req)
	if err != nil {
		log.Fatal(err)
	}

	// decode, _ := json.Marshal(res.Body)

	fmt.Println(status)

	// client := binance.NewClient(API_KEY, SECRET_KEY)

	// permission, err := client.NewAveragePriceService().Do(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// encoded, _ := json.Marshal(permission)

	// fmt.Println(string(encoded))
}
