package k8s

import "github.com/portworx/torpedo/drivers/scheduler"

// SchedName is the name of the kubernetes scheduler driver implementation
const SchedName = "k8s"

type k8s struct{}

func (k *k8s) GetNodes() ([]string, error) {
	return []string{}, nil
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
