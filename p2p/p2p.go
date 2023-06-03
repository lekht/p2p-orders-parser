package p2p

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var tradeOperation []string = []string{"buy", "sell"}

type P2PBinance struct {
	// todo давай просто стандартный клиент юзать, из стандартной либы. А логику попыток лучше реализрвать на уровне P2PBinance
	client *MyHttpClient
}

type MyHttpClient struct {
	attempts int
	http.Client
}

func NewP2PBinance() *P2PBinance {
	c := http.Client{Timeout: 10 * time.Second}
	client := MyHttpClient{attempts: 3, Client: c}

	return &P2PBinance{client: &client}
}

func (p *P2PBinance) GetOrderBooks(ctx context.Context, fiats, assets []string) (map[string]map[string]OrderBook, error) {
	book := make(map[string]map[string]OrderBook)

	var wg sync.WaitGroup

	var rp RequestParameters

	rp.Countries = nil
	rp.Page = 1
	rp.PayTypes = nil
	rp.ProMerchantAds = false
	rp.PublisherType = nil
	rp.Rows = 10
	ch := make(chan OrderBook)

	for _, asset := range assets {
		rp.Asset = asset
		for _, fiat := range fiats {
			rp.Fiat = fiat
			wg.Add(1)
			go doRequest(&wg, ch, p, rp)
		}
	}

	// todo этот wg не нужен, т/к ниже мы не в гоуьине будем читать
	var wg1 sync.WaitGroup

	wg1.Add(1)
	// todo эта горутина не должна быть горутиной. Она цикл по каналу закончится и так, когда закроется канал
	go func(map[string]map[string]OrderBook) {
		for b := range ch {
			f := b.Buy[0].Fiat
			a := b.Buy[0].Asset

			if book[f] == nil {
				book[f] = make(map[string]OrderBook)
			}

			book[f][a] = b
		}
		wg1.Done()
	}(book)

	// todo это должно перехать наверх (перед проходом по циклу и быть в горутине) либо в отедбльной функции вызываются все методы doRequest мы их ждем и потом закрываем канал
	wg.Wait()
	close(ch)

	wg1.Wait()
	return book, nil
}

// todo это должно быть методом P2PBinance
// todo над doRequest должен быть метод который и занимается бизнес логикой
// todo doRequest должен делать ровно одну вещи и делать ее хорошо, т/е сделать http запрос и вернуть ответ, распарсенный в структуру
// todo вместо клиента лучше завернуть логику попыток во враппер над этим метдом
// todo внутри не должно быть for _, operation := range tradeOperation { - кажд запрос в отдельной рутине
func doRequest(wg *sync.WaitGroup, ch chan OrderBook, p *P2PBinance, r RequestParameters) {
	var b OrderBook
	for _, operation := range tradeOperation { // todo это уже бизнес логика, а не логика запроса
		r.TradeType = operation
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(r)

		// request to binance
		req, err := http.NewRequest(http.MethodPost, "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", payloadBuf)
		if err != nil {
			// todo нигде кроме маин не должно быть паник
			log.Panicf("doRequest() - NewRequest() error: %s\n", err)
		}

		req.Header.Add("content-type", "application/json")

		response, err := p.client.do(req)
		if err != nil {
			log.Panicf("doRequest() - do() error: %s\n", err)
		}

		orders, err := dataToOrders(response)
		if err != nil {
			log.Panicf("doRequest() - dataToOrders() error: %s\n", err)
		}

		// todo это уже бизнес логика, а не логика запроса
		if r.TradeType == "buy" {
			b.Buy = orders
		} else {
			b.Sell = orders
		}
	}

	ch <- b
	wg.Done()
}

func (c *MyHttpClient) do(r *http.Request) (*Response, error) {
	var response Response
	for i := 0; i <= c.attempts; i++ {
		resp, err := c.Client.Do(r)
		if err != nil {
			continue
		}
		if resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Panic(err)
		}
		return &response, nil
	}

	return nil, errors.New("do() - all atemptes failed")
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
