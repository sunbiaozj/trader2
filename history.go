package trader

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

//LoadHistory of instrument
func (e *Engine) LoadHistory(symbol string, tf Timeframe, reader io.Reader) error {

	const format = "2006-01-02 15:04:05"

	r := csv.NewReader(reader)

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

	type tmpTF struct {
		timeframe Timeframe
		duration  time.Duration
	}

	timeframes := []tmpTF{
		tmpTF{M1, time.Minute},
		tmpTF{M5, time.Minute * 5},
		tmpTF{M15, time.Minute * 15},
		tmpTF{M30, time.Minute * 30},
		tmpTF{H1, time.Hour},
		tmpTF{H4, time.Hour * 4},
		tmpTF{D1, time.Hour * 24},
	}

	for i := 0; i < (len(timeframes) - 1); i++ {

		currDuration, currTimeFrame := timeframes[i].duration, timeframes[i].timeframe
		nextDuration, nextTimeFrame := timeframes[i+1].duration, timeframes[i+1].timeframe

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
