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
	GetProfitMatches([]matcher.FiatPairOrder) []matcher.TradeChain
}

func main() {
	// getting parameters for request

	parametersPath := flag.String("c", "conf.yml", "path to config file")
	flag.Parse()

	if *parametersPath != "" {
		err := params.Load(*parametersPath)
		if err != nil {
			// todo лучше не писать "я не смог открыть файл" "произошла ошибка" - это засоряет лог ошибки и все. Почитать статьи как грамотно писать сообщение  об ошибке
			log.Panicf("load conf: %s\n", err)
		}
	} else {
		log.Println("there are no parameters' path")
		log.Panicf("you should use the flag ---> --parameters=")
	}

	log.Println(params)

	// todo работать с интерфесом, а не структорой
	p := p2p.NewP2PBinance()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	// todo сделать вывод не руками а метод типа PrintResults
	fmt.Print("\n\n")
	fmt.Print("buy ", "\n", book["RUB"]["BTC"].Buy, "\n\n")
	fmt.Print("buy ", "\n", book["RUB"]["USDT"].Buy, "\n\n")
	fmt.Print("sell ", "\n", book["RUB"]["USDT"].Sell, "\n\n")
	fmt.Print("sell ", "\n", book["KZT"]["BTC"].Sell, "\n\n")
	fmt.Print("buy ", "\n", book["KZT"]["USDT"].Buy, "\n\n")

	//m := matcher.NewMatcher()

}

func printBook(map[string]map[string]p2p.OrderBook) {
	// todo достать из мапы
	// todo осортировать по фиату, ассету
	// todo напечатать все это в цикле
	// todo печать только верхние
	// todo печатать числа типа 1+3 как 1000
}
