package matcher

import (
	"p2p-orders-parser/p2p"
	"reflect"
	"testing"
)

func TestPriceMatcherSimple_GetFiatOrders(t *testing.T) {
	type args struct {
		m map[string]map[string]p2p.OrderBook
	}
	tests := []struct {
		name string
		p    *PriceMatcherSimple
		args args
		want []FiatPairOrder
	}{
		{
			name: "empty book",
			p:    NewMatcher(),
			args: args{
				m: map[string]map[string]p2p.OrderBook{},
			},
			want: nil,
		},
		{"1", NewMatcher(), args{map[string]map[string]p2p.OrderBook{
			"RUB": {
				"USDT": {
					Buy: []p2p.Order{
						{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU1"},
					},
					Sell: []p2p.Order{
						{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU6"},
					},
				},
				"BTC": {
					Buy: []p2p.Order{
						{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB1"},
					},
					Sell: []p2p.Order{
						{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB6"},
					},
				},
			},
			"KZT": {
				"USDT": {
					Buy: []p2p.Order{
						{Asset: "USDT", Fiat: "KZT", Price: 450.00, PaymentMethod: "", Advertiser: "KU1"},
					},
					Sell: []p2p.Order{
						{Asset: "USDT", Fiat: "KZT", Price: 460.00, PaymentMethod: "", Advertiser: "KU6"},
					},
				},
				"BTC": {
					Buy: []p2p.Order{
						{Asset: "BTC", Fiat: "KZT", Price: 13870000.00, PaymentMethod: "", Advertiser: "KB1"},
					},
					Sell: []p2p.Order{
						{Asset: "BTC", Fiat: "KZT", Price: 14000000.00, PaymentMethod: "", Advertiser: "KB6"},
					},
				},
			},
		}}, []FiatPairOrder{
			{
				Buy:  p2p.Order{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU1"},
				Sell: p2p.Order{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU6"},
			},
			{
				Buy:  p2p.Order{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB1"},
				Sell: p2p.Order{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB6"},
			},
			{
				Buy:  p2p.Order{Asset: "USDT", Fiat: "KZT", Price: 450.00, PaymentMethod: "", Advertiser: "KU1"},
				Sell: p2p.Order{Asset: "USDT", Fiat: "KZT", Price: 460.00, PaymentMethod: "", Advertiser: "KU6"},
			},
			{
				Buy:  p2p.Order{Asset: "BTC", Fiat: "KZT", Price: 13870000.00, PaymentMethod: "", Advertiser: "KB1"},
				Sell: p2p.Order{Asset: "BTC", Fiat: "KZT", Price: 14000000.00, PaymentMethod: "", Advertiser: "KB6"},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PriceMatcherSimple{}
			if got := p.GetFiatOrders(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PriceMatcherSimple.GetFiatOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceMatcherSimple_GetProfitMatches(t *testing.T) {
	type args struct {
		pairs []FiatPairOrder
	}
	tests := []struct {
		name string
		p    *PriceMatcherSimple
		args args
		want []TradeChain
	}{
		{"empty", NewMatcher(), args{pairs: []FiatPairOrder{}}, []TradeChain{}},
		{"1", NewMatcher(), args{
			pairs: []FiatPairOrder{
				{
					Buy:  p2p.Order{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU1"},
					Sell: p2p.Order{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU6"},
				},
				{
					Buy:  p2p.Order{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB1"},
					Sell: p2p.Order{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB6"},
				},
				{
					Buy:  p2p.Order{Asset: "USDT", Fiat: "KZT", Price: 450.00, PaymentMethod: "", Advertiser: "KU1"},
					Sell: p2p.Order{Asset: "USDT", Fiat: "KZT", Price: 460.00, PaymentMethod: "", Advertiser: "KU6"},
				},
				{
					Buy:  p2p.Order{Asset: "BTC", Fiat: "KZT", Price: 13870000.00, PaymentMethod: "", Advertiser: "KB1"},
					Sell: p2p.Order{Asset: "BTC", Fiat: "KZT", Price: 14000000.00, PaymentMethod: "", Advertiser: "KB6"},
				},
			},
		}, []TradeChain{
			{[]FiatPairOrder{
				{
					Buy:  p2p.Order{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU1"},
					Sell: p2p.Order{Asset: "USDT", Fiat: "KZT", Price: 460.00, PaymentMethod: "", Advertiser: "KU6"},
				},
				{
					Buy:  p2p.Order{Asset: "USDT", Fiat: "KZT", Price: 450.00, PaymentMethod: "", Advertiser: "KU1"},
					Sell: p2p.Order{Asset: "USDT", Fiat: "RUB", Price: 80.0, PaymentMethod: "", Advertiser: "RU6"},
				},
			}},
			{[]FiatPairOrder{
				{
					Buy:  p2p.Order{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB1"},
					Sell: p2p.Order{Asset: "BTC", Fiat: "RUB", Price: 2000.00, PaymentMethod: "", Advertiser: "RB6"},
				},
			}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PriceMatcherSimple{}
			if got := p.GetProfitMatches(tt.args.pairs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PriceMatcherSimple.GetProfitMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
