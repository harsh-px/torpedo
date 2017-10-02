package torpedo

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTorpedo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo testsuites")
}

var _ = BeforeSuite(func() {
	err := Instance().S.Init(Instance().SpecDir)
	Expect(err).NotTo(HaveOccurred())

	err = Instance().V.Init(Instance().S.String())
	Expect(err).NotTo(HaveOccurred())

	err = Instance().N.Init(Instance().S.String())
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// teardown code for Torpedo testsuite
})
