package matcher

import (
	"math"
	"p2p-orders-parser/config"
	"p2p-orders-parser/p2p"
)

type TradeChain struct {
	Orders []FiatPairOrder // [buy:rub-usdt sell:rub-usdt] or [buy:rub-usdt sell:kzt-usdt] [buy:kzt-btc sell:rub-btc]
}

type PriceMatcherSimple struct {
	assets []string
	fiats  []string
}

type FiatPairOrder struct {
	Buy  p2p.Order
	Sell p2p.Order
}

func NewMatcher(conf config.Conf) *PriceMatcherSimple {
	return &PriceMatcherSimple{
		assets: conf.Asset,
		fiats:  conf.Fiat,
	}

}

// returns best pairs
func (p *PriceMatcherSimple) GetFiatOrders(m map[string]map[string]p2p.OrderBook) []FiatPairOrder {
	var pairs []FiatPairOrder

	for _, fiat := range p.fiats {
		for _, asset := range p.assets {
			book := m[fiat][asset]
			for i := 0; i < len(book.Buy); i++ {
				var pair FiatPairOrder
				pair.Buy = book.Buy[i]
				pair.Sell = book.Sell[i]

				if k := pair.Sell.Price / pair.Buy.Price; k < 1 {
					continue
				}

				pairs = append(pairs, pair)
			}
		}

	}

	return pairs
}

// todo: fix
// todo: []fiats ==> [fiats[0] -> ... -> fiats[0]]

// creates possible variations of transactions
func (p *PriceMatcherSimple) GetProfitMatches(pairs []FiatPairOrder) []TradeChain {
	tradeChain := make([]TradeChain, 0)

	var buyOrders []p2p.Order
	var sellOrders []p2p.Order

	for _, pair := range pairs {
		buyCopy := pair.Buy
		sellCopy := pair.Sell

		buyOrders = append(buyOrders, buyCopy)
		sellOrders = append(sellOrders, sellCopy)
	}

	for _, o := range buyOrders {

		buy := append([]p2p.Order(nil), buyOrders...)
		sell := append([]p2p.Order(nil), sellOrders...)

	inner:
		for {
			orders, ok := searchLongChain(&buy, &sell, o)
			if !ok {
				break inner
			}
			tradeChain = append(tradeChain, TradeChain{
				[]FiatPairOrder{
					{Buy: o, Sell: orders[0]},
					{Buy: orders[1], Sell: orders[2]},
				},
			})
		}

		//searchShortChain
	}

	return tradeChain
}

func searchLongChain(buy, sell *[]p2p.Order, o p2p.Order) ([]p2p.Order, bool) {
	mainFiat, currentFiat := o.Fiat, o.Fiat
	currentAsset := o.Asset

	i := 1
	ok := true

	orders := make([]p2p.Order, 0, 4)

Loop:
	for i < 4 && ok {
		switch i {
		case 1:
			for j, s := range *sell {
				if s.Fiat != mainFiat && s.Asset == currentAsset {
					orders = append(orders, s)
					currentFiat = s.Fiat
					sell = removeIndex(sell, j)
					i++
					continue Loop
				}
			}
			ok = false
		case 2:
			for j, b := range *buy {
				if b.Fiat == currentFiat {
					orders = append(orders, b)
					currentAsset = b.Asset
					buy = removeIndex(buy, j)
					i++
					continue Loop
				}
			}
			ok = false
		case 3:
			for j, s := range *sell {
				if s.Fiat == mainFiat && s.Asset == currentAsset {
					orders = append(orders, s)
					sell = removeIndex(sell, j)
					i++
					continue Loop
				}
			}
			ok = false
		}
	}
	return orders, ok
}

func removeIndex(o *[]p2p.Order, idx int) *[]p2p.Order {
	r := make([]p2p.Order, 0)
	slice := *o
	r = append(r, slice[:idx]...)
	r = append(r, slice[idx+1:]...)
	return &r
}

func (t TradeChain) Profit() float64 {
	var c float64 = 1
	for _, o := range t.Orders {
		c /= o.Buy.Price
		c *= o.Sell.Price
	}

	return math.Round(c * 100)
}

func (t TradeChain) Fiats() []string {
	var fiats []string
	for _, o := range t.Orders {
		fiats = append(fiats, o.Buy.Fiat, o.Sell.Fiat)
	}
	return fiats
}
