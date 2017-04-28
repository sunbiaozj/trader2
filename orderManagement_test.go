package trader

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrderManagement", func() {

	symbol := "BTC/USD"

	strategy := &Strategy{
		Title:  "Some",
		Code:   "qqq",
		Symbol: symbol,
		Size:   100,
		Parts:  5,
	}

	Context("Signal validity", func() {
		It("SellOpen overlaps with BuyClose", func() {
			signal := Signal{
				BuyClose: []Level{
					Level{101, 1},
					Level{102, 1},
					Level{106, 2},
					Level{104, 1},
				},
				SellOpen: []Level{
					Level{103, 1},
				},
			}

			Expect(checkSignal(signal)).Should(HaveOccurred())
		})
	})

	Context("Orders", func() {
		It("Should create new buy orders", func() {
			ex := &exchng{
				pos: []Position{
					Position{symbol, 100, 150},
				},
				ord: []Order{
					Order{
						ID:     "67",
						Symbol: symbol,
						Price:  97,
						Amount: 200,
					}, Order{
						ID:     "5",
						Symbol: symbol,
						Price:  95.5,
						Amount: 200,
					},
				},
			}

			om := newOrderManagement(strategy, ex)
			om.sendOrdersOn = false
			go om.signalLoop()

			signal := Signal{
				Strategy: strategy,
				BuyOpen: []Level{
					Level{99, 1},
					Level{98, 1},
					Level{97, 2},
					Level{96, 1},
				},
				BuyClose: []Level{
					Level{102, 3},
					Level{103, 1},
				},
			}

			Expect(checkSignal(signal)).Should(Succeed())

			om.signalCh <- signal

			openOrders := []Order{
				Order{Symbol: symbol, Price: 98, Amount: 50},
				Order{Symbol: symbol, Price: 96, Amount: 100},
				Order{Symbol: symbol, Price: 102, Amount: -150},
			}

			cancelOrders := []Order{
				Order{ID: "5"},
			}

			time.Sleep(time.Millisecond)

			Expect(om.newOrders).Should(Equal(openOrders))
			Expect(om.cancelOrders).Should(Equal(cancelOrders))
		})

		It("Should create new sell orders", func() {

			ex := &exchng{
				pos: []Position{
					Position{symbol, 100, 150},
				},
				ord: []Order{
					Order{
						Symbol: symbol,
						Price:  97,
						Amount: 250,
					},
				},
			}

			om := newOrderManagement(strategy, ex)
			om.sendOrdersOn = false
			go om.signalLoop()

			signal := Signal{
				Strategy: strategy,
				SellOpen: []Level{
					Level{105, 1},
					Level{108, 1},
					Level{102, 2},
					Level{109, 1},
				},
				SellClose: []Level{
					Level{105, 1},
					Level{108, 1},
					Level{102, 2},
					Level{109, 1},
				},
			}

			Expect(checkSignal(signal)).Should(Succeed())

			om.signalCh <- signal

			openOrders := []Order{
				Order{Symbol: symbol, Price: 102, Amount: -350},
				Order{Symbol: symbol, Price: 105, Amount: -100},
				Order{Symbol: symbol, Price: 108, Amount: -100},
				Order{Symbol: symbol, Price: 109, Amount: -100},
			}

			time.Sleep(time.Millisecond)

			Expect(om.newOrders).Should(Equal(openOrders))
		})
	})
})
