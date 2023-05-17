package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"p2p-orders-parser/config"
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
	err := params.ReqParams()
	if err != nil {
		log.Panicf("main - new request error: %s\n", err)
	}
}

func main() {
	// getting parameters for request

	log.Println(params)

	data := params

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(data)

	req, err := http.NewRequest(http.MethodPost, "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", payloadBuf)
	if err != nil {
		log.Panicf("main - new request error: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// request to binance
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("main - client: making http request error: %s\n", err)
	}

	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(response)
}
