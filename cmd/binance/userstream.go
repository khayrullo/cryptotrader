package binance

import (
	"github.com/spf13/viper"
	"log"
	"github.com/crankykernel/cryptotrader/binance"
	"encoding/json"
)

func BinanceUserStreamCommand() {
	apiKey := viper.GetString("binance.api.key")

	if apiKey == "" {
		log.Fatal("error: this command requires an api key")
	}

	restClient := binance.NewClient(&binance.RestClientAuth{
		ApiKey: apiKey,
	})

	streamClient, err := binance.OpenUserStream(restClient)
	if err != nil {
		log.Fatal("error: failed to open user stream: %v", err)
	}

	for {
		_, body, err := streamClient.Next()
		if err != nil {
			log.Fatal("error: failed to read next message: %v", err)
		}

		var rawOrderUpdate binance.OrderUpdate
		err = json.Unmarshal(body, &rawOrderUpdate)
		if err == nil {
			x := binance.OrderUpdate(rawOrderUpdate)
			log.Printf("%+v\n", rawOrderUpdate)
			log.Printf("%+v\n", x)
		}
	}
}
