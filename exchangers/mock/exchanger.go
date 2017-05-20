package mock

import (
	"math/rand"
	"strconv"

	"github.com/apex/log"
	"github.com/santacruz123/trader"
)

type exchng struct {
	pos []trader.Position
	ord []trader.Order
}

//New constructor
func New() trader.Exchanger {
	return &exchng{}
}

func (e *exchng) FetchOrders() (o []trader.Order, err error) {
	log.Info("Fetching orders")
	return
}

func (e *exchng) FetchPositions() (p []trader.Position, err error) {
	log.Info("Fetching positions")
	return
}

func (e *exchng) Orders() (o []trader.Order) {
	return e.ord
}

func (e *exchng) Positions() (p []trader.Position) {
	return e.pos
}

func (e *exchng) NewOrder(o trader.Order) (string, error) {
	o.ID = strconv.Itoa(rand.Int())

	log.WithFields(log.Fields{
		"Amount": o.Amount,
		"ID":     o.ID,
		"Price":  o.Price,
	}).Info("New Order")

	e.ord = append(e.ord, o)
	return o.ID, nil
}

func (e *exchng) CancelOrder(ID string) (ok bool, err error) {
	log.WithFields(log.Fields{
		"ID": ID,
	}).Info("Cancel order")

	return true, nil
}
