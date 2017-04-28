package trader

//Signal - signal of strategy
type Signal struct {
	Strategy                               *Strategy
	BuyOpen, BuyClose, SellOpen, SellClose []Level
}

//Level for price/size
type Level struct {
	Price float64
	Size  int
}

type byPrice []Level

func (l byPrice) Len() int {
	return len(l)
}

func (l byPrice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l byPrice) Less(i, j int) bool {
	return l[i].Price < l[j].Price
}

type byPriceDesc []Level

func (l byPriceDesc) Len() int {
	return len(l)
}

func (l byPriceDesc) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l byPriceDesc) Less(i, j int) bool {
	return l[i].Price > l[j].Price
}
