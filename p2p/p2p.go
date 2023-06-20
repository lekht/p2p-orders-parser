package p2p

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	urlSearchOrders = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	tradeTypeBuy    = "buy"
	tradeTypeSell   = "sell"
)

var (
	ErrWrongCode        = errors.New("wrong response code")
	ErrAttemptsFailed   = errors.New("all  attempts failed")
	ErrInvalidTradeType = errors.New("invalid trade type")
	ErrNoInput          = errors.New("no input parameters")
)

type P2PBinance struct {
	attempts int
	client   *http.Client
}

func NewP2PBinance() *P2PBinance {
	c := http.Client{Timeout: 10 * time.Second}

	return &P2PBinance{attempts: 3, client: &c}
}

func (p *P2PBinance) GetOrderBooks(ctx context.Context, fiats, assets []string) (map[string]map[string]OrderBook, error) {
	if len(fiats) == 0 || len(assets) == 0 {
		return nil, errors.Wrap(ErrNoInput, "GetOrderBooks()")
	}

	book := make(map[string]map[string]OrderBook)

	rps := makeRequestParameters(fiats, assets)

	chanOrders := make(chan Orders)
	chanErr := make(chan error)

	go p.newWorkers(rps, chanOrders, chanErr)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(*sync.WaitGroup) {
		for o := range chanOrders {
			f := o.Orders[0].Fiat
			a := o.Orders[0].Asset

			if book[f] == nil {
				book[f] = make(map[string]OrderBook)
			}

			b := book[f][a]
			switch o.TradeType {
			case tradeTypeBuy:
				b.Buy = o.Orders
			case tradeTypeSell:
				b.Sell = o.Orders
			default:
				chanErr <- ErrInvalidTradeType
			}
			book[f][a] = b
		}
		wg.Done()
	}(&wg)

	for err := range chanErr {
		if err != nil {
			return nil, errors.Wrap(err, "failed to get order books")
		}
	}

	wg.Wait()

	return book, nil
}

func (p *P2PBinance) newWorkers(rps []RequestParameters, chanOrders chan<- Orders, chanErr chan<- error) {
	var wg sync.WaitGroup

	for _, rp := range rps {
		wg.Add(1)
		rp := rp
		go func() {
			defer wg.Done()
			p.getOrders(rp, chanOrders, chanErr)
		}()
	}

	wg.Wait()
	close(chanOrders)
	close(chanErr)
}

func (p *P2PBinance) getOrders(rp RequestParameters, chanOrders chan<- Orders, chanErr chan<- error) {
	resp, err := p.parse(rp)
	if err != nil {
		chanErr <- errors.Wrap(err, "getOrders() - parse()")
		return
	}

	o, err := dataToOrders(resp)
	if err != nil {
		chanErr <- errors.Wrap(err, "getOrders() - dataToOrders()")
		return
	}

	switch rp.TradeType {
	case tradeTypeBuy:
		chanOrders <- Orders{Orders: o, TradeType: tradeTypeBuy}
	case tradeTypeSell:
		chanOrders <- Orders{Orders: o, TradeType: tradeTypeSell}
	default:
		chanErr <- errors.Wrap(ErrInvalidTradeType, "getOrders()")
	}
}

func makeRequestParameters(fiats, assets []string) []RequestParameters {
	var requests []RequestParameters

	for _, f := range fiats {
		for _, a := range assets {
			requests = append(requests,
				RequestParameters{
					Page:      1,
					Rows:      10,
					Asset:     a,
					Fiat:      f,
					TradeType: tradeTypeBuy,
				},
				RequestParameters{
					Page:      1,
					Rows:      10,
					Asset:     a,
					Fiat:      f,
					TradeType: tradeTypeSell,
				})
		}
	}

	return requests
}

func (p *P2PBinance) parse(rp RequestParameters) (*Response, error) {
	for i := 0; i < p.attempts; i++ {
		resp, err := p.doRequest(rp)
		if err != nil {
			continue
		}
		return resp, nil

	}
	return nil, ErrAttemptsFailed
}

// todo внутри не должно быть for _, operation := range tradeOperation { - кажд запрос в отдельной рутине
func (p *P2PBinance) doRequest(rp RequestParameters) (*Response, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(rp)

	// request to binance
	req, err := http.NewRequest(http.MethodPost, urlSearchOrders, payloadBuf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Add("content-type", "application/json")

	var response Response
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed request to server")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Wrap(ErrWrongCode, "Do()")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body into struct")
	}
	return &response, nil
}

func dataToOrders(r *Response) ([]Order, error) {
	var orders []Order
	for _, d := range r.Objects {
		var o Order
		o.Asset = d.Adv.Asset
		o.Fiat = d.Adv.Fiat
		price, err := strconv.ParseFloat(d.Adv.Price, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert price into float64")
		}
		o.Price = price
		o.PaymentMethod = d.Adv.Trade[0].TradeName
		o.Advertiser = d.Advertiser.Nick

		orders = append(orders, o)
	}
	return orders, nil
}
