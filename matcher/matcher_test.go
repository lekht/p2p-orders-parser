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
