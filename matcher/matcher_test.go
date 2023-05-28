package matcher

import (
	"context"
	"log"
	"p2p-orders-parser/config"
	"p2p-orders-parser/p2p"
	"testing"
)

var params config.Conf

func TestPriceMatcherSimple_GetFiatOrders(t *testing.T) {
	params.Asset = []string{"USDT"}
	params.Fiat = []string{"RUB"}

	p := p2p.NewP2PBinance()

	book, err := p.GetOrderBooks(context.Background(), params.Fiat, params.Asset)
	if err != nil {
		log.Panic(err)
	}

	m := NewMatcher()

	pairs := m.GetFiatOrders(book)
	log.Println(pairs)
}
