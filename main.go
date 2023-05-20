package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"p2p-orders-parser/config"
	"time"
)

type Response struct {
	Objects []Object `json:"data"`
}

type Object struct {
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

var params config.Parameters

func init() {
	parametersPath := flag.String("parameters", "", "path to config file")
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

	log.Println(params)
	var rp RequestParameters

	rp.Countries = nil
	rp.Page = 1
	rp.PayTypes = nil
	rp.ProMerchantAds = false
	rp.PublisherType = nil
	rp.Rows = 10

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
}

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
