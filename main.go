package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"p2p-orders-parser/config"
	"p2p-orders-parser/matcher"
	"p2p-orders-parser/p2p"
	"strconv"
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

	err := params.Parse(*path)
	if err != nil {
		log.Panicf("failed to parse config: %s\n", err)
	}

	var p P2P = p2p.NewP2PBinance()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	var m PriceMatcher = matcher.NewMatcher()

	pairs := m.GetFiatOrders(book)
	result := m.GetProfitMatches(pairs)

	for i, c := range result {
		fmt.Println(i)
		fmt.Println(c)
		fmt.Println(c.Profit())
		fmt.Println(c.Fiats())
	}

	printBook(book)
}

func printBook(book map[string]map[string]p2p.OrderBook) {

	for f, fiatBook := range book {
		fmt.Print("\n=========\n", f, "\n=========\n")

		for a, b := range fiatBook {
			buy := b.Buy[0]
			sell := b.Sell[0]
			fmt.Print("----------\n", a, "\n----------\n")
			fmt.Print("BUY:\n")
			fmt.Print(buy.Asset, " ", buy.Fiat, " ", strconv.FormatFloat(buy.Price, 'f', -1, 64), " ", buy.PaymentMethod, " ", buy.Advertiser, "\n")
			fmt.Print("SELL:\n")
			fmt.Print(sell.Asset, " ", sell.Fiat, " ", sell.Price, " ", sell.PaymentMethod, " ", sell.Advertiser, "\n")
		}
		fmt.Print("\n")
	}
}
