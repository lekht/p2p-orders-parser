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
		fmt.Print(i+1, ". ")
		for _, o := range c.Orders {
			fmt.Print("BUY:[", o.Buy.Fiat, " >>> ", o.Buy.Asset, " ", strconv.FormatFloat(o.Buy.Price, 'f', -1, 64), "] ")
			fmt.Print("SELL:[", o.Sell.Asset, " >>> ", o.Sell.Fiat, " ", strconv.FormatFloat(o.Sell.Price, 'f', -1, 64), "] ")
		}
		fmt.Print("PROFIT: ", strconv.FormatFloat(c.Profit()*100, 'f', 0, 64), "%", "\n")
	}

	printFullBook(book)
}

// 1. Buy:[RUB >>> USDT 80.00] SELL:[USDT >>> RUB 82.00]
// 2. Buy:[RUB >>> USDT 80.00] SELL:[] BUY:[] SELL:[USDT >>> RUB 82.00]

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
			fmt.Print(sell.Asset, " ", sell.Fiat, " ", strconv.FormatFloat(sell.Price, 'f', -1, 64), " ", sell.PaymentMethod, " ", sell.Advertiser, "\n")
		}
		fmt.Print("\n")
	}
}

func printFullBook(book map[string]map[string]p2p.OrderBook) {

	for f, fiatBook := range book {
		fmt.Print("\n=========\n", f, "\n=========\n")

		for a, b := range fiatBook {
			fmt.Print("----------\n", a, "\n----------\n")
			fmt.Print("BUY:\n")
			for _, buy := range b.Buy {
				fmt.Print(buy.Asset, " ", buy.Fiat, " ", strconv.FormatFloat(buy.Price, 'f', -1, 64), " ", buy.PaymentMethod, " ", buy.Advertiser, "\n")
			}
			fmt.Print("SELL:\n")
			for _, sell := range b.Sell {
				fmt.Print(sell.Asset, " ", sell.Fiat, " ", strconv.FormatFloat(sell.Price, 'f', -1, 64), " ", sell.PaymentMethod, " ", sell.Advertiser, "\n")
			}
		}
		fmt.Print("\n")
	}
}

// todo находить абсолютно все цепочки, не только для рубля и не только те, где 2 сделки и в обеих рубль, н-р RUB->USDT->KZT->BTC->RUB или RUB->USDT->KZT->USDT->RUB
// todo написать тесты для матчера GetProfitMatches с разными наборами буков, в том числе там должны быть наборы которые гаранитрованно дают цепочки из п выше
// todo прикрутить докерфайл для нашего приложения
// todo создать docker-compose.yml в котором поднимать наше приложение + БД mongo
// todo при запуске приложения оно должно раз в минуту по cron делать запросы и склдывать все результаты в БД - 2 табл - books & profit_chains. В кажд добавить время, когда мы делали очередной запрос
// todo поля таблицы book: created_at, fiat, asset, buy_price, buy_available, buy_method, sell_price, sell_available, sell_method. Available должен быть выражен в asset
// todo поля таблицы profit_chain: created_at, fiat, roi, массив буков и в каждом поля: fiat, asset, side (buy/sell), price, method
// todo изучить как пользоваться дебагером и показать
// todo вывод профитных матчей причесать (убрать все вот эти скобки, ник убрать, числа порядковые убрать или писать на той же строчке, указывать где бай где сел, в самом переди писать валюту в которой мы зарабатываем, сортировать все записи по валюте в которой зарабываем, прибыль писать в процентах - это называет ROI )

// todo опционально - https://pmihaylov.com/go-service-with-elk/ - изучить статью и внедрить все что там написано
