package k8s

import (
	"github.com/portworx/torpedo/drivers/scheduler"
	"github.com/portworx/torpedo/pkg/k8sutils"
	"k8s.io/client-go/pkg/api/v1"
	"github.com/portworx/torpedo/drivers/scheduler/k8s/spec/factory"
	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

// SchedName is the name of the kubernetes scheduler driver implementation
const SchedName = "k8s"

type k8s struct{}

func (k *k8s) GetNodes() ([]scheduler.Node, error) {
	nodes, err := k8sutils.GetNodes()
	if err != nil {
		return nil, err
	}

	var retNodes []scheduler.Node
	for _, n := range nodes.Items {
		var addrs []string
		for _, addr := range n.Status.Addresses {
			if addr.Type == v1.NodeExternalIP || addr.Type == v1.NodeInternalIP {
				addrs = append(addrs, addr.Address)
			}
		}

		var nodeType scheduler.NodeType
		if k8sutils.IsNodeMaster(n) {
			nodeType = scheduler.NodeTypeMaster
		} else {
			nodeType = scheduler.NodeTypeWorker
		}

		retNodes = append(retNodes, scheduler.Node{
			Name:      n.Name,
			Addresses: addrs,
			Type: nodeType,
		})
	}

	return retNodes, nil
}

// String returns the string name of this driver.
func (k *k8s) String() string {
	return SchedName
}

func (k *k8s) Init() error {
	return nil
}

func (k *k8s) Schedule(app scheduler.App) (*scheduler.Context, error) {
	// Find spec from factory
	spec, err := factory.Get(app.Key)
	if err != nil {
		return nil, err
	}

	for _, storage := range spec.Storage() {
		logrus.Infof("Deploying storage component: %#v", storage)
	}

	for _, core := range spec.Core(1, app.Name) {
		logrus.Infof("Deploying core component: %#v", core)
	}

	ctx := &scheduler.Context{
		UID: uuid.New(),
		App: app,
		// Status: TODO
		// Stdout: TODO
		// Stderr: TODO
	}

	return ctx, nil
}

func (k *k8s) WaitDone(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) Destroy(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) GetVolumes(ctx *scheduler.Context) ([]string, error) {
	return nil, nil
}

func (k *k8s) InspectVolumes(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) DeleteVolumes(ctx *scheduler.Context) error {
	return nil
}

func init() {
	k := &k8s{}
	scheduler.Register(SchedName, k)
}
