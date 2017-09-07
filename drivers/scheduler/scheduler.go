package scheduler

import (
	"github.com/portworx/torpedo/drivers"
	"github.com/portworx/torpedo/pkg/errors"
)

// NodeType identifies the type of the cluster node
type NodeType string

const (
	// NodeTypeMaster identifies a cluster node that is a master/manager
	NodeTypeMaster NodeType = "Master"
	// NodeTypeWorker identifies a cluster node that is a worker
	NodeTypeWorker NodeType = "Worker"
)

// Node encapsulates a node in the cluster
type Node struct {
	Name      string
	Addresses []string
	Type      NodeType
}

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
	Nodes []Node
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

	// String returns the string name of this driver.
	String() string

	// GetNodes returns an array of all nodes in the cluster.
	GetNodes() ([]Node, error)

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
