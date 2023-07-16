package mongodb

import (
	"context"
	"p2p-orders-parser/matcher"
	"p2p-orders-parser/p2p"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrAttemptsFailed = errors.New("connection attempts failed")
)

type MongoDB struct {
	client *mongo.Client

	database *mongo.Database
	books    *mongo.Collection
	chain    *mongo.Collection
}

type MongoConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func New() (*MongoDB, error) {

	c, err := newClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a client")
	}
	database := c.Database("p2p")

	books := database.Collection("books")
	chain := database.Collection("profit_chains")

	db := MongoDB{
		client:   c,
		database: database,
		books:    books,
		chain:    chain,
	}

	return &db, nil
}

func newClient() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017/p2p")
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "mongo client initializing")
	}
	err = c.Ping(context.Background(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "checking connection")
	}
	return c, nil
}

func (m *MongoDB) AddBooks(b map[string]map[string]p2p.OrderBook) error {
	docs := convertBooks(b)

	_, err := m.books.InsertMany(context.Background(), docs)
	if err != nil {
		return errors.Wrap(err, "insert into collection books")
	}

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

func (m *MongoDB) AddChains(chains []matcher.TradeChain) error {
	docs := convertChains(chains)

	_, err := m.chain.InsertMany(context.Background(), docs)
	if err != nil {
		return errors.Wrap(err, "insert into collection books")
	}

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
