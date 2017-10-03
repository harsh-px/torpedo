package tests

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/portworx/torpedo"
	"github.com/portworx/torpedo/drivers/scheduler"
)

func TestBasic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo:Basic")
}

var _ = BeforeSuite(func() {
	By(fmt.Sprintf("running basic tests under torpedo instance: %v", Instance().InstanceID))
})

// testSetupTearDown performs basic test of starting an application and destroying it (along with storage)
var _ = Describe("Setup and teardown", func() {
	var err error
	var contexts []*scheduler.Context
	taskName := fmt.Sprintf("setupteardown-%v", Instance().InstanceID)

	Context("For setting up", func() {
		It("has to schedule applications", func() {
			contexts, err = Instance().S.Schedule(taskName, scheduler.ScheduleOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts).NotTo(BeEmpty())
		})

		It("has to validate the applications", func() {
			for _, ctx := range contexts {
				err = Instance().ValidateContext(ctx)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Context("For tearing down", func() {
		It("has to destroy application and it's storage", func() {
			for _, ctx := range contexts {
				err = Instance().TearDownContext(ctx)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
