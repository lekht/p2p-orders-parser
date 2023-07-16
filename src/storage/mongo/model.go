package mongodb

import "time"

// todo поля таблицы book: created_at, fiat, asset, buy_price, buy_available, buy_method, sell_price, sell_available, sell_method. Available должен быть выражен в asset

type Book struct {
	CreatedAt time.Time
	Fiat      string
	Asset     string

	BuyPrice     float64
	BuyAvailable float64
	BuyMethod    string

	SellPrice     float64
	SellAvailable float64
	SellMethod    string
}

// todo поля таблицы profit_chain: created_at, fiat, roi, массив буков и в каждом поля: fiat, asset, side (buy/sell), price, method

type ProfitChains struct {
	CreatedAt time.Time
	Roi       float64

	Orders []ProfitOrder
}

type ProfitOrder struct {
	Fiat   string
	Asset  string
	Side   string
	Price  float64
	Method string
}
