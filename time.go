package trader

import (
	"time"
)

//Timeframe type
type Timeframe uint8

//Timeframes
const (
	M1 Timeframe = iota
	M5
	M15
	M30
	H1
	H4
	D1
)

//History for symbol
type History struct {
	TimeSeries map[Timeframe]map[time.Time]*OHLC
}

//OHLC structure
type OHLC struct {
	Open, High, Low, Close, Volume float64
}

func (e *Engine) gotQuote(symbol string, quote Quote) {
	h := e.quotes[symbol]

	if quote.Bid != 0 {
		h.Bid = quote.Bid
	}

	if quote.Ask != 0 {
		h.Ask = quote.Ask
	}

}

func (e *Engine) gotTrade(symbol string, t Trade) {
	h := e.ohlc[symbol]

	tfrs := map[Timeframe]time.Duration{
		M1:  time.Minute,
		M5:  time.Minute * 5,
		M15: time.Minute * 15,
		M30: time.Minute * 30,
		H1:  time.Hour,
		H4:  time.Hour * 4,
		D1:  time.Hour * 24,
	}

	for tf, dur := range tfrs {
		tmpTime := t.Time.Truncate(dur)
		ohlc, found := h.TimeSeries[tf][tmpTime]

		if !found {
			ohlc = &OHLC{
				Open:  t.Price,
				High:  t.Price,
				Low:   t.Price,
				Close: t.Price,
			}

			h.TimeSeries[tf][tmpTime] = ohlc
		}

		if ohlc.Open == 0 {
			ohlc.Open = t.Price
		}

		if ohlc.High < t.Price {
			ohlc.High = t.Price
		}

		if ohlc.Low > t.Price {
			ohlc.Low = t.Price
		}

		ohlc.Close = t.Price
	}
}
