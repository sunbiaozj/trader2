package trader

import (
	"time"
)

//Engine that runs strategies
type Engine struct {
	strategies []*Strategy

	tradeCh map[string]chan Trade
	quoteCh map[string]chan Quote

	ohlc   map[string]*History
	quotes map[string]Quote

	quit chan struct{}
}

//Trade type
type Trade struct {
	Price, Amount float64
	Time          time.Time
}

//Quote type
type Quote struct {
	Bid, Ask float64
}

//NewEngine constructor
func NewEngine() *Engine {
	return &Engine{
		quoteCh: make(map[string]chan Quote),
		tradeCh: make(map[string]chan Trade),
		ohlc:    make(map[string]*History),
		quit:    make(chan struct{}),
	}
}

//AddSymbol to engine
func (e *Engine) AddSymbol(symbol string, quotes chan Quote, trades chan Trade) {
	e.quoteCh[symbol] = quotes
	e.tradeCh[symbol] = trades

	h := &History{
		TimeSeries: make(map[Timeframe]map[time.Time]*OHLC),
	}

	e.ohlc[symbol] = h

	h.TimeSeries = map[Timeframe]map[time.Time]*OHLC{
		M1:  make(map[time.Time]*OHLC),
		M5:  make(map[time.Time]*OHLC),
		M15: make(map[time.Time]*OHLC),
		M30: make(map[time.Time]*OHLC),
		H1:  make(map[time.Time]*OHLC),
		H4:  make(map[time.Time]*OHLC),
		D1:  make(map[time.Time]*OHLC),
	}

}

//AddStrategy to engine
func (e *Engine) AddStrategy(strategy *Strategy) {
	e.strategies = append(e.strategies, strategy)
}

//Run engine
func (e *Engine) Run() {
	go e.loop()
}

//Stop engine
func (e *Engine) Stop() {
	close(e.quit)
}

func (e *Engine) loop() {

	// Trades loop
	for symbol, tradeCh := range e.tradeCh {
		go func(s string, ch chan Trade) {
			for {
				select {
				case <-e.quit:
					return
				case trade := <-e.tradeCh[s]:
					e.gotTrade(s, trade)
				}
			}
		}(symbol, tradeCh)
	}

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			for _, one := range e.strategies {
				go func(oneStrategy *Strategy) {
					oneStrategy.Loop(e)
				}(one)
			}
		}
	}

}
