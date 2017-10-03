package tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/portworx/torpedo"
)

func TestTorpedo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo:Init")
}

var _ = BeforeSuite(func() {
})

var _ = Describe("Init", func() {
	var err error
	Context("For initialization", func() {
		It("should init scheduler driver", func() {
			err = Instance().S.Init(Instance().SpecDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should init volume driver", func() {
			err = Instance().V.Init(Instance().S.String())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get node driver", func() {
			err = Instance().N.Init(Instance().S.String())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = AfterSuite(func() {
	// teardown code for Torpedo testsuite
})
