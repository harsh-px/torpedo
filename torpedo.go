package torpedo

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/portworx/torpedo/drivers/node"
	// import aws driver to invoke it's init
	_ "github.com/portworx/torpedo/drivers/node/aws"
	// import ssh driver to invoke it's init
	_ "github.com/portworx/torpedo/drivers/node/ssh"
	"github.com/portworx/torpedo/drivers/scheduler"
	// import k8s driver to invoke it's init
	_ "github.com/portworx/torpedo/drivers/scheduler/k8s"
	"github.com/portworx/torpedo/drivers/volume"
	// import portworx driver to invoke it's init
	_ "github.com/portworx/torpedo/drivers/volume/portworx"
	"github.com/portworx/torpedo/pkg/errors"
	"github.com/portworx/torpedo/pkg/log"
)

var instance *Torpedo
var once sync.Once

var schedulerDriver scheduler.Driver
var volumeDriver volume.Driver
var nodeDriver node.Driver
var specDir *string

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

// TODO move ValidateContext, ValidateVolumes and TearDownContext to shared closures
// https://onsi.github.io/ginkgo/#shared-example-patterns

// ValidateContext validates if the context applications are running and it's volumes get provisioned
func (t *Torpedo) ValidateContext(ctx *scheduler.Context) error {
	var err error
	if ctx.Status != 0 {
		return fmt.Errorf("exit status %v\nStdout: %v\nStderr: %v",
			ctx.Status,
			ctx.Stdout,
			ctx.Stderr,
		)
	}

	if err := t.ValidateVolumes(ctx); err != nil {
		return err
	}
	if err = t.S.WaitForRunning(ctx); err != nil {
		return err
	}
	return err
}

// ValidateVolumes validates the volume with the scheduler and volume driver
func (t *Torpedo) ValidateVolumes(ctx *scheduler.Context) error {
	if err := t.S.InspectVolumes(ctx); err != nil {
		return &errors.ErrValidateVol{
			ID:    ctx.UID,
			Cause: err.Error(),
		}
	}

	// Get all volumes with their params and ask volume driver to inspect them
	volumes, err := t.S.GetVolumeParameters(ctx)
	if err != nil {
		return &errors.ErrValidateVol{
			ID:    ctx.UID,
			Cause: err.Error(),
		}
	}

	for vol, params := range volumes {
		if err := t.V.InspectVolume(vol, params); err != nil {
			return &errors.ErrValidateVol{
				ID:    ctx.UID,
				Cause: err.Error(),
			}
		}
	}

	return nil
}

// TearDownContext destroys the context applications and storage
func (t *Torpedo) TearDownContext(ctx *scheduler.Context) error {
	var err error
	if err = t.S.Destroy(ctx); err != nil {
		return err
	}

	if err = t.S.WaitForDestroy(ctx); err != nil {
		return err
	}

	if err = t.S.DeleteVolumes(ctx); err != nil {
		return err
	}

	return err
}

// Instance returns the Torpedo singleton
func Instance() *Torpedo {
	once.Do(func() {
		instance = &Torpedo{
			InstanceID: time.Now().Format("01-02-15h04m05s"),
			S:          schedulerDriver,
			V:          volumeDriver,
			N:          nodeDriver,
			SpecDir:    *specDir,
		}

		logrus.Infof("[debug] create torpedo instance: %p", instance)
	})

	return instance
}

func init() {
	var err error
	logrus.SetLevel(logrus.InfoLevel)
	logrus.StandardLogger().Hooks.Add(log.NewHook())

	s := flag.String(schedulerCliFlag, defaultScheduler, "Name of the scheduler to us")
	n := flag.String(nodeDriverCliFlag, defaultNodeDriver, "Name of the node driver to use")
	v := flag.String(storageDriverCliFlag, defaultStorageDriver, "Name of the storage driver to use")
	specDir = flag.String(specDirCliFlag, DefaultSpecsRoot,
		"Root directory container the application spec files")

	flag.Parse()

	if schedulerDriver, err = scheduler.Get(*s); err != nil {
		logrus.Fatalf("Cannot find scheduler driver for %v. Err: %v\n", *s, err)
		os.Exit(-1)
	} else if volumeDriver, err = volume.Get(*v); err != nil {
		logrus.Fatalf("Cannot find volume driver for %v. Err: %v\n", *v, err)
		os.Exit(-1)
	} else if nodeDriver, err = node.Get(*n); err != nil {
		logrus.Fatalf("Cannot find node driver for %v. Err: %v\n", *n, err)
		os.Exit(-1)
	}
}
