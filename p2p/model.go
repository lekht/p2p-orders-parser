package p2p

type Response struct {
	Objects []Data `json:"data"`
}

type Data struct {
	Adv        Adv        `json:"adv"`
	Advertiser Advertiser `json:"advertiser"`
}

type Adv struct {
	Asset     string         `json:"asset"`
	Price     string         `json:"price"`
	Fiat      string         `json:"fiatUnit"`
	Available string         `json:"surplusAmount"`
	AmountMin string         `json:"minSingleTransAmount"`
	AmountMax string         `json:"maxSingleTransAmount"`
	Trade     []TradeMethods `json:"tradeMethods"`
}

type TradeMethods struct {
	TradeName      string `json:"tradeMethodName"`
	TradeShortName string `json:"tradeMethodShortName"`
}

type Advertiser struct {
	Nick string `json:"nickName"`
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

type OrderBook struct {
	Buy  []Order
	Sell []Order
}

type Orders struct {
	Orders    []Order
	TradeType string
}

type Order struct {
	Asset         string
	Fiat          string
	Price         float64
	PaymentMethod string
	Advertiser    string
}
