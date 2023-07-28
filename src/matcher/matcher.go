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

func sortPairsByProfit(pairs []FiatPairOrder) []FiatPairOrder {
	newPairs := make([]FiatPairOrder, len(pairs))

	copy(newPairs, pairs)

	sort(newPairs, 0, len(pairs)-1)

	return newPairs
}

func sort(pairs []FiatPairOrder, start, end int) {
	if (end - start) < 1 {
		return
	}

	pivot := pairs[end].Sell.Price / pairs[end].Buy.Price
	pivotIdx := end
	splitIndex := start

	for i := start; i < end; i++ {
		if k := pairs[i].Sell.Price / pairs[i].Sell.Price; k < pivot {
			temp := pairs[splitIndex]

			pairs[splitIndex] = pairs[i]
			pairs[i] = temp

			splitIndex++
		}
	}

	pairs[end] = pairs[splitIndex]
	pairs[splitIndex] = pairs[pivotIdx]

	sort(pairs, start, splitIndex-1)
	sort(pairs, splitIndex+1, end)
}

// creates possible variations of transactions
func (p *PriceMatcherSimple) GetProfitMatches(pairs []FiatPairOrder) []TradeChain {
	tradeChain := make([]TradeChain, 0)

	var buyOrders []p2p.Order
	var sellOrders []p2p.Order

	sortedPairs := sortPairsByProfit(pairs)

	for _, pair := range sortedPairs {
		buyOrders = append(buyOrders, pair.Buy)
		sellOrders = append(sellOrders, pair.Sell)
	}

	for i, b := range buyOrders {

		buy := append([]p2p.Order(nil), buyOrders...)
		sell := append([]p2p.Order(nil), sellOrders...)

		buy[i] = p2p.Order{}

	inner:
		for {
			orders, ok := searchLongChain(buy, sell, b)
			if !ok {
				break inner
			}

			tradeChain = append(tradeChain, TradeChain{
				[]FiatPairOrder{
					{Buy: b, Sell: sellOrders[orders[0]]},
					{Buy: buyOrders[orders[1]], Sell: sellOrders[orders[2]]},
				},
			})
		}

		newSell := append([]p2p.Order(nil), sellOrders...)
		chains := searchShortChain(b, newSell)
		for _, c := range chains {
			tradeChain = append(tradeChain, TradeChain{[]FiatPairOrder{c}})
		}

	}

	return tradeChain
}

func searchShortChain(b p2p.Order, sell []p2p.Order) []FiatPairOrder {
	short := make([]FiatPairOrder, 0)

	for _, s := range sell {
		if s.Asset == b.Asset && s.Fiat == b.Fiat {
			short = append(short, FiatPairOrder{Buy: b, Sell: s})
		}
	}

	return short
}

func searchLongChain(buy, sell []p2p.Order, buyOrder p2p.Order) ([]int, bool) {
	mainFiat, currentFiat := buyOrder.Fiat, buyOrder.Fiat
	currentAsset := buyOrder.Asset

	i := 1
	ok := true

	orders := make([]int, 0, 4)

Loop:
	for i < 4 && ok {
		switch i {
		case 1:
			for j, s := range sell {
				if s.Fiat != mainFiat && s.Asset == currentAsset {
					orders = append(orders, j)
					currentFiat = s.Fiat
					sell[j] = p2p.Order{}
					i++
					continue Loop
				}
			}
			ok = false
		case 2:
			for j, b := range buy {
				if b.Fiat == currentFiat {
					orders = append(orders, j)
					currentAsset = b.Asset
					buy[j] = p2p.Order{}
					i++
					continue Loop
				}
			}
			ok = false
		case 3:
			for j, s := range sell {
				if s.Fiat == mainFiat && s.Asset == currentAsset {
					orders = append(orders, j)
					sell[j] = p2p.Order{}

					i++
					continue Loop
				}
			}
			ok = false
		}
	}
	return orders, ok
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
