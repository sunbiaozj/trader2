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
