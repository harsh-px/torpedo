package torpedo

import (
	"flag"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/portworx/torpedo/drivers/node"
	// blank importing drivers so they get inited
	_ "github.com/portworx/torpedo/drivers/node/aws"
	_ "github.com/portworx/torpedo/drivers/node/ssh"
	"github.com/portworx/torpedo/drivers/scheduler"
	// blank importing drivers so they get inited
	_ "github.com/portworx/torpedo/drivers/scheduler/k8s"
	"github.com/portworx/torpedo/drivers/volume"
	// blank importing drivers so they get inited
	_ "github.com/portworx/torpedo/drivers/volume/portworx"
	"github.com/portworx/torpedo/pkg/log"
)

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

var instance Torpedo

// Torpedo is the torpedo testsuite
type Torpedo struct {
	InstanceID string
	S          scheduler.Driver
	V          volume.Driver
	N          node.Driver
	SpecDir    string
}

// Instance returns the Torpedo singleton
func Instance() Torpedo {
	return instance
}

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.StandardLogger().Hooks.Add(log.NewHook())

	s := flag.String(schedulerCliFlag, defaultScheduler, "Name of the scheduler to us")
	n := flag.String(nodeDriverCliFlag, defaultNodeDriver, "Name of the node driver to use")
	v := flag.String(storageDriverCliFlag, defaultStorageDriver, "Name of the storage driver to use")
	specDir := flag.String(specDirCliFlag, DefaultSpecsRoot, "Root directory container the application spec files")

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
			InstanceID: time.Now().Format("01-02-15h04m05s"),
			S:          schedulerDriver,
			V:          volumeDriver,
			N:          nodeDriver,
			SpecDir:    *specDir,
		}
	}

	logrus.Printf("Torpedo initialized with volume driver: %v, and scheduler: %v, node: %v\n",
		Instance().V.String(),
		Instance().S.String(),
		Instance().N.String(),
	)
}
