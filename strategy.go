package trader

//Strategy struct
type Strategy struct {
	Symbol  string
	Code    string
	Title   string
	Parts   int
	Size    int
	Decimal int

	OnTick func(*Engine) (Signal, error)
	Init   func(*Engine) error
}
