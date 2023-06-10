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

type Parser interface {
	GetOrderBooks(ctx context.Context, fiats, assets []string) (map[string]map[string]*p2p.OrderBook, error) // fiat->asset->[]Order
}

type Matcher interface {
	GetFiatOrders(map[string]map[string]p2p.OrderBook) []matcher.FiatPairOrder
	GetProfitMatches([]matcher.FiatPairOrder) []matcher.TradeChain
}

func main() {
	// getting parameters for request
	parametersPath := flag.String("c", "conf.yml", "path to config file")
	flag.Parse()

	if *parametersPath == "" {
		log.Println("there are no parameters' path")
		log.Panicf("you should use the flag ---> --parameters=")
	}

	err := params.ReqParams(*parametersPath)
	if err != nil {
		log.Panicf("main - new request error: %s\n", err)
	}

	log.Println(params)

	if err := newApplication().Run(); err != nil {
		log.Panic(err)
	}
}

type application struct {
	parser  Parser
	matcher Matcher
}

func newApplication() *application {
	a := &application{}

	//init dependencies
	a.initParser()
	a.initMatcher()

	return a
}

func (a *application) initParser() {
	a.parser = p2p.NewP2PBinance()
}

func (a *application) initMatcher() {
	a.matcher = matcher.NewMatcher()
}

func (a *application) Run() error {
	book, err := a.parser.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		return err
	}

	fmt.Print("\n\n")
	fmt.Print("buy ", "\n", book["RUB"]["USDT"].Buy, "\n\n")
	fmt.Print("sell ", "\n", book["RUB"]["USDT"].Sell, "\n\n")
	fmt.Print("buy ", "\n", book["KZT"]["USDT"].Buy, "\n\n")
	fmt.Print("buy ", "\n", book["RUB"]["BTC"].Buy, "\n\n")
	fmt.Print("sell ", "\n", book["KZT"]["BTC"].Sell, "\n\n")

	// m := matcher.NewMatcher()

	return nil
}

// todo create constructors for all services
