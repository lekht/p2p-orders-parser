package storage

import (
	"p2p-orders-parser/matcher"
	"p2p-orders-parser/p2p"
	dummydb "p2p-orders-parser/storage/dummy"

	"github.com/pkg/errors"
)

type DB interface {
	AddBooks(b map[string]map[string]p2p.OrderBook) error
	AddChains(chains []matcher.TradeChain) error
}

type Storage struct {
	db DB
}

func New() (*Storage, error) {
	db, err := dummydb.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create niw storage")
	}

	return &Storage{db: db}, nil
}

func (s Storage) AddBooks(b map[string]map[string]p2p.OrderBook) error {
	return s.db.AddBooks(b)
}

func (s Storage) AddChains(chains []matcher.TradeChain) error {
	return s.db.AddChains(chains)
}
