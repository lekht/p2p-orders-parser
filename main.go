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
	path := flag.String("c", "", "path to config file")
	flag.Parse()

	if *path != "" {
		err := params.Load(*path)
		if err != nil {
			// todo лучше не писать "я не смог открыть файл" "произошла ошибка" - это засоряет лог ошибки и все. Почитать статьи как грамотно писать сообщение  об ошибке
			log.Panicln("Load: ", err)
		}
	} else {
		log.Println("there are no parameters' path")
		log.Panicf("you should use the flag ---> --parameters=")
	}

	log.Println(params)

	var p P2P = p2p.NewP2PBinance()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	var m PriceMatcher = matcher.NewMatcher()

	pairs := m.GetFiatOrders(book)
	_ = m.GetProfitMatches(pairs)

	printBook(book)

}

func printBook(book map[string]map[string]p2p.OrderBook) {
	// todo достать из мапы
	// todo осортировать по фиату, ассету
	// todo напечатать все это в цикле
	// todo печать только верхние
	// todo печатать числа типа 1+3 как 1000

	for f, fiat_book := range book {
		fmt.Print("\n=========\n", f, "\n=========\n")

		for a, asset := range fiat_book {
			fmt.Print("----------\n", a, "\n----------\n")
			fmt.Print("BUY:\n", asset.Buy[0], "\nSELL:\n", asset.Sell[0], "\n")
		}
		fmt.Print("\n")
	}
}
