package scheduler

import (
	"github.com/portworx/torpedo/drivers"
	"github.com/portworx/torpedo/pkg/errors"
)

const (
	// LocalHost will pin a task to the node the task is created on.
	LocalHost = "localhost"
	// ExternalHost will pick any other host in the cluster other than the
	// one the task is created on.
	ExternalHost = "externalhost"
)

// Volume specifies the parameters for creating an external volume.
type Volume struct {
	Driver string
	Name   string
	Size   int // in GB
}

// App encapsulates an application run within a scheduler
type App struct {
	ID       string
	Name     string
	Replicas int
	Vol      Volume
	// Nodes in which to run the task. If empty, scheduler will pick the node(s).
	Nodes []string
}

// Context holds the execution context and output values of a test task.
type Context struct {
	ID     string
	App    App
	Status int
	Stdout string
	Stderr string
}

// Driver must be implemented to provide test support to various schedulers.
type Driver interface {
	// Driver provides the basic service manipulation routines.
	drivers.Driver

	// GetNodes returns an array of all nodes in the cluster.
	GetNodes() ([]string, error)

	// Create creates a task context. Does not start the task.
	Create(App) (*Context, error)

	// Schedule starts a task
	Schedule(*Context) error

	// WaitDone waits for task to complete.
	WaitDone(*Context) error

	// Destroy removes a task. Must also delete the external volume.
	Destroy(*Context) error

	// InspectVolume inspects a storage volume.
	InspectVolume(name string) error

	// DeleteVolume will delete a storage volume.
	DeleteVolume(name string) error
}

var (
	schedulers = make(map[string]Driver)
)

// Register registers the given scheduler driver
func Register(name string, d Driver) error {
	schedulers[name] = d
	return nil
}

// Get returns a registered scheduler test provider.
func Get(name string) (Driver, error) {
	if d, ok := schedulers[name]; ok {
		return d, nil
	}
	return nil, &errors.ErrNotFound{
		ID:   name,
		Type: "Scheduler",
	}
}
