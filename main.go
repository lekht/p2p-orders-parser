package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"p2p-orders-parser/config"
	"p2p-orders-parser/matcher"
	"p2p-orders-parser/p2p"
)

var params config.Conf

type P2P interface {
	GetOrderBooks(ctx context.Context, fiats, assets []string) (map[string]map[string]p2p.OrderBook, error) // fiat->asset->[]Order
}

type PriceMatcher interface {
	GetFiatOrders(map[string]map[string]p2p.OrderBook) []matcher.FiatPairOrder
	GetProfitMatches([]matcher.FiatPairOrder) []matcher.FiatPairOrderProfit
}

func main() {
	// getting parameters for request

	parametersPath := flag.String("c", "conf.yml", "path to config file")
	flag.Parse()

	if *parametersPath != "" {
		err := params.ReqParams(*parametersPath)
		if err != nil {
			log.Panicf("main - new request error: %s\n", err)
		}
	} else {
		log.Println("there are no parameters' path")
		log.Panicf("you should use the flag ---> --parameters=")
	}

	log.Println(params)

	p := p2p.NewP2PBinance()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}
	fmt.Print("\n\n")
	fmt.Print("buy ", "\n", book["RUB"]["USDT"].Buy, "\n\n")
	fmt.Print("sell ", "\n", book["RUB"]["USDT"].Sell, "\n\n")
	fmt.Print("buy ", "\n", book["KZT"]["USDT"].Buy, "\n\n")
	fmt.Print("buy ", "\n", book["RUB"]["BTC"].Buy, "\n\n")
	fmt.Print("sell ", "\n", book["KZT"]["BTC"].Sell, "\n\n")

}

// todo create constructors for all services
