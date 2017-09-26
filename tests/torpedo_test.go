package tests

import (
	"flag"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/portworx/torpedo/drivers/node"
	"github.com/portworx/torpedo/drivers/scheduler"
	"github.com/portworx/torpedo/drivers/volume"
)

const (
	schedulerCliFlag     = "scheduler"
	nodeDriverCliFlag    = "node-driver,n"
	storageDriverCliFlag = "storage,v"
	testsCliFlag         = "tests,t"
)

const (
	defaultScheduler     = "k8s"
	defaultNodeDriver    = "ssh"
	defaultStorageDriver = "pxd"
)

var instance Torpedo

type Torpedo struct {
	instanceID string
	s          scheduler.Driver
	v          volume.Driver
	n          node.Driver
}

func TestTorpedo(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Torpedo testsuite")
}

var _ = ginkgo.BeforeSuite(func() {
	// setup code for torpedo testsuite
})

var _ = ginkgo.AfterSuite(func() {
	// teardown code for torpedo testsuite
})

// Instance returns the torpedo singleton
func Instance() Torpedo {
	return instance
}

func init() {
	s := flag.String(schedulerCliFlag, defaultScheduler, "Name of the scheduler to us")
	n := flag.String(nodeDriverCliFlag, defaultNodeDriver, "Name of the node driver to use")
	v := flag.String(storageDriverCliFlag, defaultStorageDriver, "Name of the storage driver to use")

	flag.Parse()

	if schedulerDriver, err := scheduler.Get(*s); err != nil {
		logrus.Fatalf("Cannot find scheduler driver for %v. Err: %v\n", *s, err)
		os.Exit(-1)
	} else if volumeDriver, err := volume.Get(*v); err != nil {
		logrus.Fatalf("Cannot find volume driver for %v. Err: %v\n", *v, err)
		os.Exit(-1)
	} else if nodeDriver, err := node.Get(*n); err != nil {
		logrus.Fatalf("Cannot find node driver for %v. Err: %v\n", *n, err)
		os.Exit(-1)
	} else {
		instance = Torpedo{
			instanceID: time.Now().Format("01-02-15h04m05s"),
			s:          schedulerDriver,
			v:          volumeDriver,
			n:          nodeDriver,
		}
	}

	logrus.Printf("Torpedo initialized with volume driver: %v, and scheduler: %v, node: %v\n",
		Instance().v.String(),
		Instance().s.String(),
		Instance().n.String(),
	)
}
