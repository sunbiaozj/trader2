package trader

//Strategy - holds
type Strategy struct {
	Symbol string
	Code   string
	Title  string
	Sizes  uint
	OnTick func(*Engine) (Signal, error)
}
