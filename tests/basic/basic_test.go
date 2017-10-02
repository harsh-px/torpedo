package tests

import (
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
	//"github.com/Sirupsen/logrus"
)

func TestBasic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo:Basic")
}

var _ = BeforeSuite(func() {
	//logrus.Infof("running basic tests under torpedo instance: %v", Instance().InstanceID)
})

var _ = Describe("Setup and teardown", func() {
	Context("initially", func() {
		It("has 0 items", func() {})
		It("has 0 units", func() {})
		Specify("the total amount is 0.00", func() {})
	})
})

