package k8sutils

import (
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	storage_v1beta1 "k8s.io/client-go/pkg/apis/storage/v1beta1"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

const k8sMasterLabelKey = "node-role.kubernetes.io/master"
const k8sPVCStorageClassKey = "volume.beta.kubernetes.io/storage-class"

// GetK8sClient instantiates a k8s client
func GetK8sClient() (*kubernetes.Clientset, error) {
	k8sClient, err := loadClientFromServiceAccount()
	if err != nil {
		return nil, err
	}

	if k8sClient == nil {
		return nil, ErrK8SApiAccountNotSet
	}

	return k8sClient, nil
}

// GetNodes talks to the k8s api server and gets the nodes in the cluster
func GetNodes() (*v1.NodeList, error) {
	var err error
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	nodes, err := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// CreateDeployment creates the given deployment
func CreateDeployment(deployment *v1beta1.Deployment) (*v1beta1.Deployment, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return client.AppsV1beta1().Deployments(deployment.Namespace).Create(deployment)
}

// DeleteDeployment deletes the given deployment
func DeleteDeployment(deployment *v1beta1.Deployment) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	return client.AppsV1beta1().Deployments(deployment.Namespace).Delete(deployment.Name, &meta_v1.DeleteOptions{})
}

// ValidateDeployement validates the given deployment if it's running and healthy
func ValidateDeployement(deployment *v1beta1.Deployment) error {
	task := func() error {
		client, err := GetK8sClient()
		if err != nil {
			return err
		}

		dep, err := client.AppsV1beta1().Deployments(deployment.Namespace).Get(deployment.Name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		if *dep.Spec.Replicas == dep.Status.AvailableReplicas {
			return &ErrAppNotReady{
				ID:    dep.Name,
				Cause: fmt.Sprintf("Expected replicas: %v Available replicas: %v", dep.Spec.Replicas, dep.Status.AvailableReplicas),
			}
		}

		if *dep.Spec.Replicas == dep.Status.ReadyReplicas {
			return &ErrAppNotReady{
				ID:    dep.Name,
				Cause: fmt.Sprintf("Expected replicas: %v Ready replicas: %v", dep.Spec.Replicas, dep.Status.ReadyReplicas),
			}
		}

		// TODO perform deeper checks with pods
		return nil
	}

	if err := doRetryWithTimeout(task, 1*time.Minute, 10*time.Second); err != nil {
		return err
	}

	return nil
}

// CreateStorageClass creates the given storage class
func CreateStorageClass(sc *storage_v1beta1.StorageClass) (*storage_v1beta1.StorageClass, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return client.StorageV1beta1().StorageClasses().Create(sc)
}

// DeleteStorageClass deletes the given storage class
func DeleteStorageClass(sc *storage_v1beta1.StorageClass) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	return client.StorageV1beta1().StorageClasses().Delete(sc.Name, &meta_v1.DeleteOptions{})
}

// CreatePersistentVolumeClaim creates the given persistent volume claim
func CreatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return client.PersistentVolumeClaims(pvc.Namespace).Create(pvc)
}

// DeletePersistentVolumeClaim deletes the given persistent volume claim
func DeletePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	return client.PersistentVolumeClaims(pvc.Namespace).Delete(pvc.Name, &meta_v1.DeleteOptions{})
}

// ValidatePersistentVolumeClaim validates the given pvc
func ValidatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) error {
	task := func() error {
		client, err := GetK8sClient()
		if err != nil {
			return err
		}

		result, err := client.PersistentVolumeClaims(pvc.Namespace).Get(pvc.Name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		if result.Status.Phase == v1.ClaimBound {
			return nil
		}

		return &ErrPVCNotReady{
			ID:    result.Name,
			Cause: fmt.Sprintf("PVC expected status: %v PVC actual status: %v", v1.ClaimBound, result.Status.Phase),
		}
	}

	if err := doRetryWithTimeout(task, 1*time.Minute, 10*time.Second); err != nil {
		return err
	}

	return nil
}

// GetVolumeForPersistentVolumeClaim returns the back volume for the given PVC
func GetVolumeForPersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (string, error) {
	client, err := GetK8sClient()
	if err != nil {
		return "", err
	}

	result, err := client.PersistentVolumeClaims(pvc.Namespace).Get(pvc.Name, meta_v1.GetOptions{})
	if err != nil {
		return "", err
	}

	return result.Spec.VolumeName, nil
}

// GetPersistentVolumeClaimParams fetches custom parameters for the given PVC
func GetPersistentVolumeClaimParams(pvc *v1.PersistentVolumeClaim) (map[string]string, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	result, err := client.PersistentVolumeClaims(pvc.Namespace).Get(pvc.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var params map[string]string

	storageResource, ok := result.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	if !ok {
		return nil, fmt.Errorf("failed to get storage resource for pvc: %v", result.Name)
	}

	params["size"] = storageResource.String()
	scName, ok := result.Annotations[k8sPVCStorageClassKey]
	if !ok {
		return nil, fmt.Errorf("failed to get storage class for pvc: %v", result.Name)
	}

	sc, err := client.StorageV1beta1().StorageClasses().Get(scName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for key, value := range sc.Parameters {
		params[key] = value
	}

	return params, nil
}

// IsNodeMaster returns true if given node is a kubernetes master node
func IsNodeMaster(node v1.Node) bool {
	_, ok := node.Labels[k8sMasterLabelKey]
	return ok
}

// loadClientFromServiceAccount loads a k8s client from a ServiceAccount specified in the pod running px
func loadClientFromServiceAccount() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}

func doRetryWithTimeout(task func() error, timeout, timeBeforeRetry time.Duration) error {
	done := make(chan bool, 1)
	quit := make(chan bool, 1)

	go func(done, quit chan bool) {
		for {
			select {
			case q := <-quit:
				if q {
					log.Printf("Quiting task due to timeout...\n")
					return
				}

			default:
				if err := task(); err == nil {
					log.Printf("Task succeeded.\n")
					done <- true
				} else {
					log.Printf("Task failing with err: %v\n", err)
				}

				time.Sleep(timeBeforeRetry)
			}
		}
	}(done, quit)

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		quit <- true
		return ErrTimedOut
	}
}
