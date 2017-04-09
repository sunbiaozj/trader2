package trader

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("History", func() {
	Context("CSV", func() {
		symbol := "BTC/USD"

		It("Loading", func() {

			ngn := NewEngine()
			quotes := make(chan Quote)
			trades := make(chan Trade)
			ngn.AddSymbol(symbol, quotes, trades)

			f, err := os.Open("fixtures/m30.csv")
			Expect(err).To(Succeed())

			err = ngn.LoadHistory(symbol, M30, f)
			Expect(err).To(Succeed())

			t, err := time.Parse("2006-01-02 15:04:05", "2017-04-07 16:30:00")

			Expect(ngn.ohlc[symbol].TimeSeries[M30][t].Open).Should(BeNumerically("==", 1193.99))
			Expect(ngn.ohlc[symbol].TimeSeries[M30][t].High).Should(BeNumerically("==", 1193.99))
			Expect(ngn.ohlc[symbol].TimeSeries[M30][t].Low).Should(BeNumerically("==", 1190.00))
			Expect(ngn.ohlc[symbol].TimeSeries[M30][t].Close).Should(BeNumerically("==", 1191.94))
			Expect(ngn.ohlc[symbol].TimeSeries[M30][t].Volume).Should(BeNumerically("==", 93))
		})

		It("Calculate", func() {

			ngn := NewEngine()
			quotes := make(chan Quote)
			trades := make(chan Trade)
			ngn.AddSymbol(symbol, quotes, trades)

			f, err := os.Open("fixtures/m30.csv")
			Expect(err).To(Succeed())

			err = ngn.LoadHistory(symbol, M30, f)
			Expect(err).To(Succeed())

			t, err := time.Parse("2006-01-02 15:04:05", "2017-04-07 04:00:00")

			// 2017-04-07 04:00:00,1186.13,1194.98,1184.87,1194.72,40
			// 2017-04-07 04:30:00,1193.83,1195.87,1189.00,1190.74,44
			// 2017-04-07 05:00:00,1190.31,1194.93,1190.31,1194.92,75
			// 2017-04-07 05:30:00,1194.93,1197.10,1188.15,1190.33,98
			// 2017-04-07 06:00:00,1191.37,1191.37,1185.24,1186.12,71
			// 2017-04-07 06:30:00,1185.02,1187.56,1176.65,1177.25,178
			// 2017-04-07 07:00:00,1177.43,1180.67,1176.00,1180.67,51
			// 2017-04-07 07:30:00,1181.37,1184.96,1181.37,1184.91,58

			Expect(ngn.ohlc[symbol].TimeSeries[H4][t].Open).Should(BeNumerically("==", 1186.13))
			Expect(ngn.ohlc[symbol].TimeSeries[H4][t].High).Should(BeNumerically("==", 1197.10))
			Expect(ngn.ohlc[symbol].TimeSeries[H4][t].Low).Should(BeNumerically("==", 1176.00))
			Expect(ngn.ohlc[symbol].TimeSeries[H4][t].Close).Should(BeNumerically("==", 1184.91))
			Expect(ngn.ohlc[symbol].TimeSeries[H4][t].Volume).Should(BeNumerically("==", 615))

		})
	})
})
