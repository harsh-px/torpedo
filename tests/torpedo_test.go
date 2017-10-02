package tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//. "github.com/portworx/torpedo"
	"flag"
	"github.com/portworx/torpedo/drivers/scheduler"
	"github.com/portworx/torpedo/drivers/volume"
	"github.com/portworx/torpedo/drivers/node"
	"github.com/portworx/torpedo/pkg/log"
	"github.com/Sirupsen/logrus"
)


var instance Torpedo

// Torpedo is the torpedo testsuite
type Torpedo struct {
	InstanceID string
	S          scheduler.Driver
	V          volume.Driver
	N          node.Driver
	SpecDir    string
}

const (
	// DefaultSpecsRoot specifies the default location of the base specs directory in the Torpedo container
	DefaultSpecsRoot     = "/specs"
	schedulerCliFlag     = "scheduler"
	nodeDriverCliFlag    = "node-driver,n"
	storageDriverCliFlag = "storage,v"
	specDirCliFlag       = "spec-dir"
)

const (
	defaultScheduler     = "k8s"
	defaultNodeDriver    = "ssh"
	defaultStorageDriver = "pxd"
)

func TestTorpedo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torpedo:Init")
}

var _ = BeforeSuite(func() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.StandardLogger().Hooks.Add(log.NewHook())

/*	err := Instance().S.Init(Instance().SpecDir)
	Expect(err).NotTo(HaveOccurred())

	err = Instance().V.Init(Instance().S.String())
	Expect(err).NotTo(HaveOccurred())

	err = Instance().N.Init(Instance().S.String())
	Expect(err).NotTo(HaveOccurred())*/
})

// Instance returns the Torpedo singleton
func Instance() Torpedo {
	return instance
}


var _ = Describe("Init", func() {
	var err error
	instance = Torpedo{}
	var s, n, v *string

	Context("Initially", func() {
		It("should parse command-line flags", func() {
			s = flag.String(schedulerCliFlag, defaultScheduler, "Name of the scheduler to us")
			n = flag.String(nodeDriverCliFlag, defaultNodeDriver, "Name of the node driver to use")
			v = flag.String(storageDriverCliFlag, defaultStorageDriver, "Name of the storage driver to use")
			instance.SpecDir = *flag.String(specDirCliFlag, DefaultSpecsRoot,
				"Root directory container the application spec files")

			flag.Parse()
		})
	})

	Context("After parsing command-line flags", func() {
		It("should get scheduler driver", func() {
			instance.S, err = scheduler.Get(*s)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get volume driver", func() {
			instance.V, err = volume.Get(*v)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get node driver", func() {
			instance.N, err = node.Get(*n)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("After getting the drivers", func() {
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


func init() {}