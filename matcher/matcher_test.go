package matcher

import (
	"context"
	"fmt"
	"log"
	"p2p-orders-parser/config"
	"p2p-orders-parser/p2p"
	"testing"
)

var params config.Conf

func initTestsDependencies() (*p2p.P2PBinance, *PriceMatcherSimple) {
	params.Fiat = []string{"RUB", "KZT", "USD"}
	params.Asset = []string{"USDT", "BTC"}
	p := p2p.NewP2PBinance()
	m := NewMatcher()

	return p, m
}

func TestPriceMatcherSimple_GetFiatOrders(t *testing.T) {

	p, m := initTestsDependencies()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	pairs := m.GetFiatOrders(book)
	log.Println("\n", pairs)

}

func TestPriceMatcherSimple_GetProfitMatches(t *testing.T) {

	p, m := initTestsDependencies()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	pairs := m.GetFiatOrders(book)

	chain := m.GetProfitMatches(pairs)
	for i, c := range chain {
		fmt.Println(i)
		fmt.Println(c)
	}
}

func TestTradeChain_Profit(t *testing.T) {
	p, m := initTestsDependencies()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	pairs := m.GetFiatOrders(book)

	chain := m.GetProfitMatches(pairs)
	for _, c := range chain {
		fmt.Println(c.Profit())
		fmt.Println(c.Fiats())
	}
}
