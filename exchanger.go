package trader

//Order abstration
type Order struct {
	ID     string
	Symbol string
	Price  float64
	Amount int
	Type   OrderType
}

//Position abstraction
type Position struct {
	Symbol string
	Price  float64
	Amount int
}

//OrderType for orders
type OrderType uint

//Limit,Stop,Market types
const (
	Limit OrderType = iota
	Stop
	Market
)

//Exchanger interface
type Exchanger interface {
	FetchOrders() ([]Order, error)
	FetchPositions() ([]Position, error)

	Orders() []Order
	Positions() []Position

	NewOrder(Order) (string, error)
	CancelOrder(string) (bool, error)
}

type exchng struct {
	pos []Position
	ord []Order
}

func (e *exchng) FetchOrders() (o []Order, err error) {
	return
}

func (e *exchng) FetchPositions() (p []Position, err error) {
	return
}

func (e *exchng) Orders() (o []Order) {
	return e.ord
}

func (e *exchng) Positions() (p []Position) {
	return e.pos
}

func (e *exchng) NewOrder(o Order) (ID string, err error) {
	return
}

func (e *exchng) CancelOrder(ID string) (ok bool, err error) {
	return
}
