package main

import (
	"fmt"
	"log"

	"github.com/izacgaldino23/binance-consult-trade-api/config"
	"github.com/izacgaldino23/binance-consult-trade-api/utils"
)

const (
	API_KEY    = "5288aVerH00ezxl7bvtB0bdNqfqnBkDuL2OnZs6E42NijebgeLuDXwdSKEiU8NCk"
	SECRET_KEY = "nxMiP7rUa7MmfFnlJcRcx74nw0YSiFBLas3G3gvTYOrfEx5SpxuqvRk8J9uMoPh9"
)

func main() {
	// LOAD ENVIRONMENT
	err := config.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	_, status, err := utils.Get("/v3/ping", &utils.ParamList{}, &utils.ParamList{})
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
