package p2p

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	urlSearchOrders = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"

	tradeTypeBuy  = "buy"
	tradeTypeSell = "sell"
)

var (
	ErrAllAttemptsFail  = errors.New("all attempts failed")
	ErrInvalidTradeType = errors.New("invalid trade type")
)

type P2PBinance struct {
	attempts   int
	httpClient *http.Client
}

func NewP2PBinance() *P2PBinance {
	c := http.Client{Timeout: 10 * time.Second}

	return &P2PBinance{
		attempts:   3,
		httpClient: &c,
	}
}

func (p *P2PBinance) GetOrderBooks(ctx context.Context, fiats, assets []string) (map[string]map[string]*OrderBook, error) {
	book := make(map[string]map[string]*OrderBook)

	var wg sync.WaitGroup

	chanOrders := make(chan Orders)

	reqParams := p.makeRequestParameters(fiats, assets)

	for _, rp := range reqParams {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.getOrders(chanOrders, rp)
		}()
	}

	for orders := range chanOrders {
		f := orders.Orders[0].Fiat
		a := orders.Orders[0].Asset

		if book[f] == nil {
			book[f] = make(map[string]*OrderBook)
			book[f][a] = &OrderBook{}
		}

		b := book[f][a]
		switch orders.TradeType {
		case tradeTypeBuy:
			b.Buy = orders.Orders
		case tradeTypeSell:
			b.Buy = orders.Orders
		default:
			return nil, ErrInvalidTradeType
		}
	}

	wg.Wait()
	close(chanOrders)

	return book, nil
}

func (p *P2PBinance) makeRequestParameters(fiats, assets []string) []RequestParameters {
	var reqParams []RequestParameters
	for _, asset := range assets {
		for _, fiat := range fiats {
			reqParams = append(reqParams, RequestParameters{
				Page:      1,
				Rows:      10,
				Asset:     asset,
				Fiat:      fiat,
				TradeType: "buy",
			},
				RequestParameters{
					Page:      1,
					Rows:      10,
					Asset:     asset,
					Fiat:      fiat,
					TradeType: "sell",
				})
		}
	}
	return reqParams
}

func (p *P2PBinance) getOrders(ch chan<- Orders, r RequestParameters) error {
	o, err := p.doRequest(r)
	if err != nil {
		return err
	}

	var orders = Orders{
		Orders: o,
	}

	switch r.TradeType {
	case tradeTypeBuy:
		orders.TradeType = tradeTypeBuy
	case tradeTypeSell:
		orders.TradeType = tradeTypeSell
	default:
		return ErrInvalidTradeType
	}

	ch <- orders
	return nil
}

func (p *P2PBinance) doRequest(r RequestParameters) ([]Order, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(r)

	// doRequest to binance
	req, err := http.NewRequest(http.MethodPost, urlSearchOrders, payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("doRequest() - NewRequest() error: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	response, err := p.do(req)
	if err != nil {
		return nil, fmt.Errorf("doRequest() - do() error: %w", err)
	}

	orders, err := dataToOrders(response)
	if err != nil {
		return nil, fmt.Errorf("doRequest() - dataToOrders() error: %w", err)
	}

	return orders, nil
}

func (p *P2PBinance) do(r *http.Request) (*Response, error) {
	var response Response
	for i := 0; i < p.attempts; i++ {
		resp, err := p.httpClient.Do(r)
		if err != nil {
			continue
		}
		if resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}

	return nil, fmt.Errorf("do() - %w", ErrAllAttemptsFail)
}

func dataToOrders(r *Response) ([]Order, error) {
	var orders []Order
	for _, d := range r.Objects {
		var o Order
		o.Asset = d.Adv.Asset
		o.Fiat = d.Adv.Fiat
		price, err := strconv.ParseFloat(d.Adv.Price, 64)
		if err != nil {
			return nil, errors.New("dataToOrders() error: price string convert error")
		}
		o.Price = price
		o.PaymentMethod = d.Adv.Trade[0].TradeName
		o.Advertiser = d.Advertiser.Nick

		orders = append(orders, o)
	}
	return orders, nil
}
