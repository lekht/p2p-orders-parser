package dummydb

import (
	"fmt"
	"p2p-orders-parser/matcher"
	"p2p-orders-parser/p2p"
	"time"
)

type DummyStorage struct {
}

func New() (*DummyStorage, error) {
	return &DummyStorage{}, nil
}

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

func (d *DummyStorage) AddBooks(b map[string]map[string]p2p.OrderBook) error {
	fmt.Println(convertBooks(b)...)

	return nil
}

func convertBooks(orderBooks map[string]map[string]p2p.OrderBook) []interface{} {
	books := []interface{}{}

	for _, fb := range orderBooks {
		for _, ab := range fb {
			for i := 0; i < len(ab.Buy); i++ {
				books = append(books, Book{
					CreatedAt: time.Now(),
					Fiat:      ab.Buy[i].Fiat,
					Asset:     ab.Buy[i].Asset,

					BuyPrice:     ab.Buy[i].Price,
					BuyAvailable: ab.Buy[i].Available,
					BuyMethod:    ab.Buy[i].PaymentMethod,

					SellPrice:     ab.Sell[i].Price,
					SellAvailable: ab.Sell[i].Available,
					SellMethod:    ab.Sell[i].PaymentMethod,
				})
			}
		}

	}

	return books
}

func (d *DummyStorage) AddChains(chains []matcher.TradeChain) error {
	fmt.Println(convertChains(chains)...)

	return nil

}

func convertChains(chains []matcher.TradeChain) []interface{} {
	profitChains := []interface{}{}

	for _, c := range chains {
		var orders []ProfitOrder

		for _, o := range c.Orders {
			orders = append(orders,
				ProfitOrder{
					Fiat:   o.Buy.Fiat,
					Asset:  o.Buy.Asset,
					Side:   "buy",
					Price:  o.Buy.Price,
					Method: o.Buy.PaymentMethod,
				},
				ProfitOrder{
					Fiat:   o.Sell.Fiat,
					Asset:  o.Sell.Asset,
					Side:   "sell",
					Price:  o.Sell.Price,
					Method: o.Sell.PaymentMethod,
				},
			)
		}

		profitChains = append(profitChains, ProfitChains{
			CreatedAt: time.Now(),
			Roi:       c.Profit(),

			Orders: orders,
		})
	}

	return profitChains
}
