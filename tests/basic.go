package tests

import (
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("Setup and teardown", func() {
	ginkgo.Context("initially", func() {
		ginkgo.It("has 0 items", func() {})
		ginkgo.It("has 0 units", func() {})
		ginkgo.Specify("the total amount is 0.00", func() {})
	})
})
