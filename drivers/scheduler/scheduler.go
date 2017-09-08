package scheduler

import (
	"log"

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

// App encapsulates an application run within a scheduler
type App struct {
	Key  string
	Name string
	// Nodes in which to run the task. If empty, scheduler will pick the node(s).
	Nodes []Node
}

// Context holds the execution context and output values of a test task.
type Context struct {
	UID    string
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
	GetNodes() []Node

	// Schedule starts a task
	Schedule(App) (*Context, error)

	// WaitForRunning waits for task to complete. TODO add WaitOptions and use vendored px sched pkg to implement
	WaitForRunning(*Context) error

	// Destroy removes a task. It does not delete the volumes of the task.
	Destroy(*Context) error

	// Returns list of volume IDs using by given context
	GetVolumes(*Context) ([]string, error)

	// InspectVolumes inspects a storage volume.
	InspectVolumes(*Context) error

	// DeleteVolumes will delete a storage volume.
	DeleteVolumes(*Context) error
}

var (
	schedulers = make(map[string]Driver)
)

// Register registers the given scheduler driver
func Register(name string, d Driver) error {
	log.Printf("Registering sched driver: %v\n", name)
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
