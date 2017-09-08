package k8sutils

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	storage_v1beta1 "k8s.io/client-go/pkg/apis/storage/v1beta1"
	"k8s.io/client-go/rest"
)

const k8sMasterLabelKey = "node-role.kubernetes.io/master"

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
