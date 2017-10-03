package tests

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/portworx/torpedo"
)

var context = ginkgo.Context
var it = ginkgo.It
var expect = gomega.Expect
var haveOccurred = gomega.HaveOccurred
var instance = torpedo.Instance

/*func TestTorpedo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo:Init")
}*/

/*var _ = BeforeSuite(func() {})

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

		It("should init node driver", func() {
			err = Instance().N.Init(Instance().S.String())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})*/

// InitInstance is the ginkgo spec for initializing torpedo
func InitInstance() func() {
	return func() {
		var err error
		context("For initialization", func() {
			it("should init scheduler driver", func() {
				err = instance().S.Init(instance().SpecDir)
				expect(err).NotTo(haveOccurred())
			})

			it("should init volume driver", func() {
				err = instance().V.Init(instance().S.String())
				expect(err).NotTo(haveOccurred())
			})

			it("should init node driver", func() {
				err = instance().N.Init(instance().S.String())
				expect(err).NotTo(haveOccurred())
			})
		})
	}
}

/*var _ = AfterSuite(func() {
	// teardown code for Torpedo testsuite
})*/
