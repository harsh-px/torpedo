package k8s

import (
	"github.com/portworx/torpedo/drivers/scheduler"
	"github.com/portworx/torpedo/pkg/k8sutils"
	"k8s.io/client-go/pkg/api/v1"
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

func (k *k8s) Create(app scheduler.App) (*scheduler.Context, error) {
	return nil, nil
}

func (k *k8s) Schedule(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) WaitDone(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) Destroy(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) InspectVolume(name string) error {
	return nil
}

func (k *k8s) DeleteVolume(name string) error {
	return nil
}

func init() {
	k := &k8s{}
	scheduler.Register(SchedName, k)
}
