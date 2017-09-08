package k8s

import (
	"fmt"
	"log"

	"github.com/pborman/uuid"
	"github.com/portworx/torpedo/drivers/scheduler"
	"github.com/portworx/torpedo/drivers/scheduler/k8s/spec/factory"
	"github.com/portworx/torpedo/pkg/k8sutils"
	"k8s.io/client-go/pkg/api/v1"
	storage_v1beta1 "k8s.io/client-go/pkg/apis/storage/v1beta1"
	// blank importing all applications specs to allow them to init()
	_ "github.com/portworx/torpedo/drivers/scheduler/k8s/spec/postgres"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

// SchedName is the name of the kubernetes scheduler driver implementation
const SchedName = "k8s"

type k8s struct {
	nodes []scheduler.Node
}

func (k *k8s) GetNodes() []scheduler.Node {
	return k.nodes
}

// String returns the string name of this driver.
func (k *k8s) String() string {
	return SchedName
}

func (k *k8s) Init() error {
	nodes, err := k8sutils.GetNodes()
	if err != nil {
		return err
	}

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

		k.nodes = append(k.nodes, scheduler.Node{
			Name:      n.Name,
			Addresses: addrs,
			Type:      nodeType,
		})
	}

	return nil
}

func (k *k8s) Schedule(app scheduler.App) (*scheduler.Context, error) {
	spec, err := factory.Get(app.Key)
	if err != nil {
		return nil, err
	}

	for _, storage := range spec.Storage() {
		if obj, ok := storage.(*storage_v1beta1.StorageClass); ok {
			sc, err := k8sutils.CreateStorageClass(obj)
			if err != nil {
				return nil, &ErrFailedToScheduleApp{
					App:   app,
					Cause: fmt.Sprintf("Failed to create storage class: %v. Err: %v", obj.Name, err),
				}
			}
			log.Printf("Created storage class: %v", sc)
		} else if obj, ok := storage.(*v1.PersistentVolumeClaim); ok {
			pvc, err := k8sutils.CreatePersistentVolumeClaim(obj)
			if err != nil {
				return nil, &ErrFailedToScheduleApp{
					App:   app,
					Cause: fmt.Sprintf("Failed to create PVC: %v. Err: %v", obj.Name, err),
				}
			}
			log.Printf("Created PVC: %v", pvc)
		} else {
			return nil, &ErrFailedToScheduleApp{
				App:   app,
				Cause: fmt.Sprintf("Failed to create unsupported storage component: %#v.", storage),
			}
		}
	}

	for _, core := range spec.Core(1, app.Name) {
		if obj, ok := core.(*v1beta1.Deployment); ok {
			dep, err := k8sutils.CreateDeployment(obj)
			if err != nil {
				return nil, &ErrFailedToScheduleApp{
					App:   app,
					Cause: fmt.Sprintf("Failed to create Deployment: %v. Err: %v", obj.Name, err),
				}
			}
			log.Printf("Created deployment: %v", dep)
		} else {
			return nil, &ErrFailedToScheduleApp{
				App:   app,
				Cause: fmt.Sprintf("Failed to create unsupported core component: %#v.", core),
			}
		}
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
	spec, err := factory.Get(ctx.App.Key)
	if err != nil {
		return err
	}

	for _, core := range spec.Core(1, ctx.App.Name) {
		if obj, ok := core.(*v1beta1.Deployment); ok {
			err := k8sutils.DeleteDeployment(obj)
			if err != nil {
				return &ErrFailedToDestroyApp{
					App:   ctx.App,
					Cause: fmt.Sprintf("Failed to destroy Deployment: %v. Err: %v", obj.Name, err),
				}
			}
			log.Printf("Destroyed deployment: %v", obj.Name)
		} else {
			return &ErrFailedToDestroyApp{
				App:   ctx.App,
				Cause: fmt.Sprintf("Failed to destroy unsupported core component: %#v.", core),
			}
		}
	}

	return nil
}

func (k *k8s) GetVolumes(ctx *scheduler.Context) ([]string, error) {
	return nil, nil
}

func (k *k8s) InspectVolumes(ctx *scheduler.Context) error {
	return nil
}

func (k *k8s) DeleteVolumes(ctx *scheduler.Context) error {
	spec, err := factory.Get(ctx.App.Key)
	if err != nil {
		return err
	}

	for _, storage := range spec.Storage() {
		if obj, ok := storage.(*storage_v1beta1.StorageClass); ok {
			if err := k8sutils.DeleteStorageClass(obj); err != nil {
				return &ErrFailedToDestroyStorage{
					App:   ctx.App,
					Cause: fmt.Sprintf("Failed to destroy storage class: %v. Err: %v", obj.Name, err),
				}
			}
			log.Printf("Destroyed storage class: %v", obj.Name)
		} else if obj, ok := storage.(*v1.PersistentVolumeClaim); ok {
			if err := k8sutils.DeletePersistentVolumeClaim(obj); err != nil {
				return &ErrFailedToDestroyStorage{
					App:   ctx.App,
					Cause: fmt.Sprintf("Failed to destroy PVC: %v. Err: %v", obj.Name, err),
				}
			}
			log.Printf("Destroyed PVC: %v", obj.Name)
		} else {
			return &ErrFailedToDestroyStorage{
				App:   ctx.App,
				Cause: fmt.Sprintf("Failed to destroy unsupported storage component: %#v.", storage),
			}
		}
	}

	return nil
}

func init() {
	k := &k8s{}
	scheduler.Register(SchedName, k)
}
