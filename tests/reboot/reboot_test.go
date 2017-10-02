package tests

import (
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
	. "github.com/portworx/torpedo"
	"github.com/Sirupsen/logrus"
)

func TestReboot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo:Reboot")
}

var _ = BeforeSuite(func() {
	logrus.Infof("running reboot tests under torpedo instance: %v", Instance().InstanceID)
})

var _ = Describe("Setup and teardown", func() {
	Context("initially", func() {
		It("has 0 items", func() {})
		It("has 0 units", func() {})
		Specify("the total amount is 0.00", func() {})
	})
})


