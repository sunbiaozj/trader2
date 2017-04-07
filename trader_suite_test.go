package trader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTrader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trader Suite")
}
