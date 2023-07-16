package p2p

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestP2PBinance_GetOrderBooks(t *testing.T) {
	type fields struct {
		attempts int
		client   *http.Client
	}
	type args struct {
		ctx    context.Context
		fiats  []string
		assets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]map[string]OrderBook
		wantErr bool
	}{
		{
			name:    "empty input",
			fields:  fields{3, &http.Client{Timeout: time.Second * 3}},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "no fiats and assets",
			fields: fields{3, &http.Client{Timeout: time.Second * 3}},
			args: args{
				ctx:    context.Background(),
				fiats:  []string{},
				assets: []string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "only fiats",
			fields: fields{3, &http.Client{Timeout: time.Second * 3}},
			args: args{
				ctx:    context.Background(),
				fiats:  []string{"RUB", "KZT"},
				assets: []string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "only assets",
			fields: fields{3, &http.Client{Timeout: time.Second * 3}},
			args: args{
				ctx:    context.Background(),
				fiats:  []string{},
				assets: []string{"USDT", "BTC"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &P2PBinance{
				attempts: tt.fields.attempts,
				client:   tt.fields.client,
			}
			got, err := p.GetOrderBooks(tt.args.ctx, tt.args.fiats, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("P2PBinance.GetOrderBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("P2PBinance.GetOrderBooks() = %v, want %v", got, tt.want)
			}
		})
	}
}
