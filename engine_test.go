package trader

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Engine", func() {

	Context("Symbol", func() {
		It("Add", func() {
			ngn := NewEngine()
			quotes := make(chan Quote)
			trades := make(chan Trade)
			ngn.AddSymbol("EUR/USD", quotes, trades)
		})

		It("OHLC", func() {

			symbol := "EUR/USD"

			ngn := NewEngine()
			quotes := make(chan Quote)
			trades := make(chan Trade)
			ngn.AddSymbol(symbol, quotes, trades)
			ngn.Run()
			defer ngn.Stop()

			t := time.Now()

			trade := Trade{
				Price: 54.23,
				Time:  t,
			}

			trades <- trade

			tM1 := t.Truncate(time.Minute)
			Expect(ngn.ohlc[symbol].TimeSeries[M1][tM1].Open).Should(BeNumerically("==", trade.Price))

			tH1 := t.Truncate(time.Hour)
			Expect(ngn.ohlc[symbol].TimeSeries[H1][tH1].Open).Should(BeNumerically("==", trade.Price))
		})
	})

	It("Quotes", func() {

		symbol := "EUR/USD"

		ngn := NewEngine()
		quotes := make(chan Quote)
		trades := make(chan Trade)
		ngn.AddSymbol(symbol, quotes, trades)
		ngn.Run()
		defer ngn.Stop()

		quote := Quote{
			Bid: 35.23,
			Ask: 36.21,
		}

		quotes <- quote

		Expect(ngn.quotes[symbol].Bid).Should(BeNumerically("==", quote.Bid))
		Expect(ngn.quotes[symbol].Ask).Should(BeNumerically("==", quote.Ask))
	})
})
