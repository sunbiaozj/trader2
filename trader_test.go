package trader_test

import (
	"github.com/santacruz123/trader"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Trader", func() {

	Context("Create", func() {
		It("Simple", func() {

			ma := &trader.Strategy{
				Title: "Moving average",
				Code:  "ma",
				Size:  3,
			}

			ngn := trader.NewEngine()
			ngn.AddStrategy(ma)
			ngn.Run()
			defer ngn.Stop()
		})
	})
})
