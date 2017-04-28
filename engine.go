package trader

import (
	"time"

	"github.com/apex/log"
)

//Engine that runs strategies
type Engine struct {
	strategy *Strategy
	om       *orderManagement

	signals chan Signal

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
	//TODO refactor into int
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
		signals:  make(chan Signal),
	}
}

//AddSymbol to engine - quotes and trades for this symbol
func (e *Engine) AddSymbol(symbol string, quotes chan Quote, trades chan Trade) {
	e.quoteCh[symbol] = quotes
	e.tradeCh[symbol] = trades
	e.changeCh[symbol] = make(chan struct{})

	initTimeSeries(e, symbol)

	e.quotes = make(map[string]*Quote)
	e.quotes[symbol] = &Quote{}
}

//AddStrategy to engine
func (e *Engine) AddStrategy(strategy *Strategy) {
	e.strategy = strategy
}

//AddExchanger to engine
func (e *Engine) AddExchanger(ex Exchanger) {
	e.om = newOrderManagement(e.strategy, ex)
	go e.om.signalLoop()
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
					if e.strategy == nil {
						continue
					}

					if e.strategy.Symbol == symbol {
						log.WithFields(log.Fields{
							"strategy": e.strategy.Code,
						}).Debug("Tick")
						go signal(e)
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
			log.WithFields(log.Fields{
				"strategy": e.strategy.Code,
			}).Debug("Ticker")
			go signal(e)
		}
	}
}

func signal(e *Engine) {
	signal, err := e.strategy.OnTick(e)

	if err != nil {
		log.WithFields(log.Fields{
			"strategy": e.strategy.Code,
		}).Errorf("Error - %v", err)
		return
	}

	select {
	case e.signals <- signal:
	default:
	}
}
