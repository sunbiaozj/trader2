package trader

import (
	"time"

	"github.com/apex/log"
)

//Engine that runs strategies
type Engine struct {
	strategies []*Strategy

	tradeCh  map[string]chan Trade
	quoteCh  map[string]chan Quote
	changeCh map[string]chan struct{}

	ohlc   map[string]*History
	quotes map[string]*Quote

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
		quoteCh:  make(map[string]chan Quote),
		tradeCh:  make(map[string]chan Trade),
		changeCh: make(map[string]chan struct{}),
		ohlc:     make(map[string]*History),
		quit:     make(chan struct{}),
	}
}

//AddSymbol to engine
func (e *Engine) AddSymbol(symbol string, quotes chan Quote, trades chan Trade) {
	e.quoteCh[symbol] = quotes
	e.tradeCh[symbol] = trades

	initTimeSeries(e, symbol)

	e.quotes = make(map[string]*Quote)
	e.quotes[symbol] = &Quote{}
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
				case trade := <-ch:
					log.WithFields(log.Fields{
						"symbol": s,
						"price":  trade.Price,
						"amount": trade.Amount,
						"time":   trade.Time,
					}).Debug("Trade")

					e.gotTrade(s, trade)
					e.changeCh[s] <- struct{}{}
				}
			}
		}(symbol, tradeCh)
	}

	// Quotes loop
	for symbol, quoteCh := range e.quoteCh {
		go func(s string, ch chan Quote) {
			for {
				select {
				case <-e.quit:
					return
				case quote := <-ch:
					log.WithFields(log.Fields{
						"symbol": s,
						"bid":    quote.Bid,
						"ask":    quote.Ask,
					}).Debug("Quote")

					e.gotQuote(s, quote)
					e.changeCh[s] <- struct{}{}
				}
			}
		}(symbol, quoteCh)
	}

	// Change
	for s := range e.tradeCh {
		go func(symbol string) {
			for {
				select {
				case <-e.quit:
					return
				case <-e.changeCh[symbol]:
					for i := range e.strategies {
						if e.strategies[i].Symbol == symbol {
							log.WithFields(log.Fields{
								"strategy": e.strategies[i],
							}).Debug("Tick")
							go e.strategies[i].OnTick(e)
						}
					}
				}
			}
		}(s)
	}

	// Run every second
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			for _, one := range e.strategies {
				go func(oneStrategy *Strategy) {
					oneStrategy.OnTick(e)
				}(one)
			}
		}
	}

}
