package trader

//Strategy - holds
type Strategy struct {
	Symbols []string
	Code    string
	Title   string
	Sizes   uint
	Loop    func(*Engine) error
}
