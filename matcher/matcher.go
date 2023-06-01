package matcher

import (
	"p2p-orders-parser/p2p"
)

type TradeChain struct {
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

	for _, fiatBooks := range m {
		for _, assetBook := range fiatBooks {
			for i := 0; i < len(assetBook.Buy); i++ {
				var pair FiatPairOrder
				pair.Buy = assetBook.Buy[i]
				pair.Sell = assetBook.Sell[i]

				if k := pair.Sell.Price / pair.Buy.Price; k < 1 {
					continue
				}

				pairs = append(pairs, pair)
			}
		}

	}

	return pairs
}

var buyOrders []p2p.Order
var sellOrders []p2p.Order

// creates possible variations of transactions
func (p *PriceMatcherSimple) GetProfitMatches(pairs []FiatPairOrder) []TradeChain {
	tradeChain := make([]TradeChain, 0)

	for _, pair := range pairs {
		buyOrders = append(buyOrders, pair.Buy)
		sellOrders = append(sellOrders, pair.Sell)

		if pair.Buy.Fiat == "RUB" {
			tradeChain = append(tradeChain, TradeChain{[]FiatPairOrder{pair}})
		}
	}

	var chainCount int

	for j := 0; j < len(buyOrders); j++ {
		if buyOrders[j].Fiat != "RUB" {
			continue
		}

		orders, ok := searchOrders(buyOrders, sellOrders, "", "", 0)
		if !ok {
			continue
		}
		var p1 FiatPairOrder = FiatPairOrder{Buy: orders[0], Sell: orders[1]}
		var p2 FiatPairOrder = FiatPairOrder{Buy: orders[2], Sell: orders[3]}
		var chain []FiatPairOrder
		chain = append(chain, p1, p2)

		tradeChain[chainCount].orders = chain
		chainCount++
	}

	return tradeChain
}

func searchOrders(buy, sell []p2p.Order, fiat, asset string, i int) ([]p2p.Order, bool) {
	var result []p2p.Order
	var ok bool

	switch i {

	case 0:
		for j, o := range buy {
			if o.Fiat == "RUB" {
				ok = true
				result = append(result, o)
				fiat = o.Fiat
				asset = o.Asset
				i++

				buyOrders = removeIndex(buyOrders, j)

				break
			}
		}

	case 1:
		// sell asset || f != fiat
		for j, o := range sell {
			if o.Asset == asset && o.Fiat != "RUB" {
				ok = true
				result = append(result, o)
				fiat = o.Fiat
				asset = o.Asset
				i++

				sellOrders = removeIndex(sellOrders, j)

				break
			}
		}

	case 2:
		// buy fiat || != asset
		for j, o := range buy {
			if o.Fiat == fiat && o.Asset != asset {
				ok = true
				result = append(result, o)
				fiat = o.Fiat
				asset = o.Asset
				i++

				buyOrders = removeIndex(buyOrders, j)

				break
			}
		}

	case 3:
		// sell asset || fiat == "RUB"
		for j, o := range sell {
			if o.Fiat == "RUB" && o.Asset == asset {
				ok = true
				result = append(result, o)
				sellOrders = removeIndex(sellOrders, j)

				return result, ok
			}
		}
	}

	if !ok {
		return nil, ok
	}

	orders, ok := searchOrders(buyOrders, sellOrders, fiat, asset, i)

	if !ok {
		return nil, ok
	}

	return append(result, orders...), ok
}

func removeIndex(s []p2p.Order, index int) []p2p.Order {
	if index == len(s)-1 {
		return s[:index]
	}
	return append(s[:index], s[index+1:]...)
}

//rub-usdt-kzt 5 buy  5000kzt -> 1000 rub
//rub-btc-kzt 6 sell  1000rub -> 6000 kzt

func (t TradeChain) Profit() float64 {
	var c float64 = 1
	for _, o := range t.orders {
		c /= o.Buy.Price
		c *= o.Sell.Price
	}
	return c
}

func (t TradeChain) Fiats() []string {
	var fiats []string
	for _, o := range t.orders {
		fiats = append(fiats, o.Buy.Fiat, o.Sell.Fiat)
	}
	return fiats
}
