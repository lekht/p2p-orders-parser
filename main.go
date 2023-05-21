package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"p2p-orders-parser/config"
	"time"
)

type Response struct {
	Objects []Data `json:"data"`
}

type Data struct {
	Adv        Adv        `json:"adv"`
	Advertiser Advertiser `json:"advertiser"`
}

type Adv struct {
	Asset     string `json:"asset"`
	Price     string `json:"price"`
	Fiat      string `json:"fiatUnit"`
	Available string `json:"surplusAmount"`
	AmountMin string `json:"minSingleTransAmount"`
	AmountMax string `json:"maxSingleTransAmount"`
}

type Advertiser struct {
	Advertiser string `json:"nickName"`
}

var params config.Conf

func init() {
	//parametersPath := flag.String("c", "conf.yml", "path to config file")
	//flag.Parse()
	//
	//if *parametersPath != "" {
	//	err := params.ReqParams(*parametersPath)
	//	if err != nil {
	//		log.Panicf("main - new request error: %s\n", err)
	//	}
	//} else {
	//	log.Println("there are no parameters' path")
	//	log.Panicf("you should use the flag ---> --parameters=")
	//}
}

type RequestParameters struct {
	Asset          string   `json:"asset"`
	Countries      []string `json:"countries"`
	Fiat           string   `json:"fiat"`
	Page           int      `json:"page"`
	PayTypes       []string `json:"payTypes"`
	ProMerchantAds bool     `json:"proMerchantAds"`
	PublisherType  []string `json:"publisherType"`
	Rows           int      `json:"rows"`
	TradeType      string   `json:""`
}

var tradeOperation []string = []string{"buy", "sell"}

func main() {
	// getting parameters for request

	parametersPath := flag.String("c", "conf.yml", "path to config file")
	flag.Parse()

	if *parametersPath != "" {
		err := params.ReqParams(*parametersPath)
		if err != nil {
			log.Panicf("main - new request error: %s\n", err)
		}
	} else {
		log.Println("there are no parameters' path")
		log.Panicf("you should use the flag ---> --parameters=")
	}

	log.Println(params)
	var rp RequestParameters

	rp.Countries = nil
	rp.Page = 1
	rp.PayTypes = nil
	rp.ProMerchantAds = false
	rp.PublisherType = nil
	rp.Rows = 10

	//ch := make(chan OrderBook)

	for _, asset := range params.Asset {
		rp.Asset = asset
		fmt.Println("==================================")
		fmt.Println(asset)
		fmt.Println("==================================")
		for _, fiat := range params.Fiat {
			rp.Fiat = fiat
			fmt.Println("--------------------------")
			fmt.Println(fiat)
			fmt.Println("--------------------------")
			for _, operation := range tradeOperation {
				rp.TradeType = operation

				fmt.Print(operation+"\n", DoRequest(&rp), "\n")
			}
		}
	}

	//for b := range ch {
	//	update map[string]map[string]OrderBook
	//}
}

// todo create constructors for all services

func DoRequest(r *RequestParameters) Response {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(r)

	c := http.Client{
		Timeout: time.Second * 10,
	}

	// request to binance
	resp, err := c.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", payloadBuf)
	if err != nil {
		log.Panicf("main - client: making http request error: %s\n", err)
	}

	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Panic(err)
	}

	return response
}

type OrderBook struct {
	Buy  []Order
	Sell []Order
}

type Order struct {
	Asset         string
	Fiat          string
	Price         float64
	Advertiser    string
	PaymentMethod string
}

type P2P interface {
	GetOrderBooks(ctx context.Context, assets, fiats []string) (map[string]map[string]OrderBook, error) // fiat->asset->[]Order
}

type P2PBinance struct {
	client MyHttpClient
}

func NewP2PBinance(client MyHttpClient) *P2PBinance {
	return &P2PBinance{client: client}
}

type MyHttpClient struct {
	attempts int
	http.Client
}

func (c MyHttpClient) Do(r *http.Request) (*http.Response, error) {
	for i := 0; i <= c.attempts; i++ {
		return c.Client.Do(r) // todo check resp code and errors
	}
	return nil, errors.New("all atemptes failed")
}

func (p P2PBinance) GetOrderBooks(ctx context.Context, assets, fiats []string) (map[string]map[string]OrderBook, error) {
	//TODO implement me
	panic("implement me")
}

type FiatPairOrderProfit struct {
	orders []FiatPairOrder // [buy:rub-usdt sell:rub-usdt] or [buy:rub-usdt sell:kzt-usdt] [buy:kzt-btc sell:rub-btc]
}

func (p FiatPairOrderProfit) Profit() float64 {
	//TODO implement me
	panic("implement me")
	for _, o := range p.orders {
		_ = o
		// todo calculate profit
	}
	return 0
}

func (p FiatPairOrderProfit) Fiats() (fiat1, fiat2 string) {
	//TODO implement me
	panic("implement me")
	return "", ""
}

type FiatPairOrder struct {
	Buy  Order
	Sell Order
}

type PriceMatcher interface {
	GetFiatOrders(map[string]map[string]OrderBook) []FiatPairOrder
	GetProfitMatches([]FiatPairOrder) []FiatPairOrderProfit
}

type PriceMatcherSimple struct{}

func (PriceMatcherSimple) GetFiatOrders(m map[string]map[string]OrderBook) []FiatPairOrder {
	//TODO implement me
	panic("implement me")
}

func (PriceMatcherSimple) GetProfitMatches(orders []FiatPairOrder) []FiatPairOrderProfit {
	//TODO implement me
	panic("implement me")
}

//rub-usdt-kzt 5 buy  5000kzt -> 1000 rub
//rub-btc-kzt 6 sell  1000rub -> 6000 kzt
