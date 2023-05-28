package matcher

import (
	"p2p-orders-parser/p2p"
)

type FiatPairOrderProfit struct {
	orders []FiatPairOrder // [buy:rub-usdt sell:rub-usdt] or [buy:rub-usdt sell:kzt-usdt] [buy:kzt-btc sell:rub-btc]
}

//rub = USDT-kzt = BTC-rub

type PriceMatcherSimple struct{}

type FiatPairOrder struct {
	Buy  p2p.Order
	Sell p2p.Order
}

// type PriceMatcher interface {
// 	GetFiatOrders(map[string]map[string]p2p.OrderBook) []matcher.FiatPairOrder
// 	GetProfitMatches([]matcher.FiatPairOrder) []matcher.FiatPairOrderProfit
// }

func NewMatcher() *PriceMatcherSimple {
	return &PriceMatcherSimple{}
}

// returns best pairs
func (p *PriceMatcherSimple) GetFiatOrders(m map[string]map[string]p2p.OrderBook) []FiatPairOrder {
	var pairs []FiatPairOrder

	return pairs
}

func (p *PriceMatcherSimple) GetProfitMatches(orders []FiatPairOrder) []FiatPairOrderProfit {
	//TODO implement me
	panic("implement me")
}

//rub-usdt-kzt 5 buy  5000kzt -> 1000 rub
//rub-btc-kzt 6 sell  1000rub -> 6000 kzt

func (p FiatPairOrderProfit) Profit() []float64 {
	//TODO implement me
	panic("implement me")
	for _, o := range p.orders {
		_ = o
		// todo calculate profit
	}
	return nil
}

func (p FiatPairOrderProfit) Fiats() (fiat1, fiat2 string) {
	//TODO implement me
	panic("implement me")
	return "", ""
}
