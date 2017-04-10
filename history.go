package trader

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

func initTimeSeries(e *Engine, s string) {
	h := &History{
		TimeSeries: make(map[Timeframe]map[time.Time]*OHLC),
	}

	e.ohlc[s] = h

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

//LoadHistory of instrument
func (e *Engine) LoadHistory(symbol string, tf Timeframe, reader io.Reader) error {

	const format = "2006-01-02 15:04:05"

	r := csv.NewReader(reader)

	initTimeSeries(e, symbol)

	for {
		rec, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		ohlc := &OHLC{}

		t, err := time.Parse(format, rec[0])

		if err != nil {
			return err
		}

		if n, err := strconv.ParseFloat(rec[1], 64); err == nil {
			ohlc.Open = n
		}

		if n, err := strconv.ParseFloat(rec[2], 64); err == nil {
			ohlc.High = n
		}

		if n, err := strconv.ParseFloat(rec[3], 64); err == nil {
			ohlc.Low = n
		}

		if n, err := strconv.ParseFloat(rec[4], 64); err == nil {
			ohlc.Close = n
		}

		if n, err := strconv.ParseFloat(rec[5], 64); err == nil {
			ohlc.Volume = n
		}

		e.ohlc[symbol].TimeSeries[tf][t] = ohlc
	}

	e.calculateTimeframes(symbol)

	return nil
}

func (e *Engine) calculateTimeframes(s string) {

	timeframes := []Timeframe{M1, M5, M15, M30, H1, H4, D1}
	durations := []time.Duration{
		time.Minute,
		time.Minute * 5,
		time.Minute * 15,
		time.Minute * 30,
		time.Hour,
		time.Hour * 4,
		time.Hour * 24,
	}

	for i := 0; i < (len(timeframes) - 1); i++ {

		currDuration, nextDuration := durations[i], durations[i+1]
		currTimeFrame, nextTimeFrame := timeframes[i], timeframes[i+1]

		currTS := e.ohlc[s].TimeSeries[currTimeFrame]
		nextTS := e.ohlc[s].TimeSeries[nextTimeFrame]

		for currTime, currOHLC := range currTS {

			nextTime := currTime.Truncate(nextDuration)
			nextOHLC, found := nextTS[nextTime]

			if !found {
				nextOHLC = &OHLC{}
				nextTS[nextTime] = nextOHLC
			}

			if nextTime == currTime {
				nextOHLC.Open = currOHLC.Open
			}

			if nextOHLC.High == 0 || nextOHLC.High < currOHLC.High {
				nextOHLC.High = currOHLC.High
			}

			if nextOHLC.Low == 0 || nextOHLC.Low > currOHLC.Low {
				nextOHLC.Low = currOHLC.Low
			}

			nextOHLC.Volume = nextOHLC.Volume + currOHLC.Volume
		}

		for nextTime, nextOHLC := range nextTS {
			subTime := nextTime.Add(nextDuration).Add(-currDuration)
			if subOHLC, found := currTS[subTime]; found {
				nextOHLC.Close = subOHLC.Close
			}
		}

	}
}
