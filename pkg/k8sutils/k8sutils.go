package k8sutils

import (
	"errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ErrK8SApiAccountNotSet is returned when the account used to talk to k8s api is not setup
var ErrK8SApiAccountNotSet = errors.New("K8s api account is not setup")

// GetK8sClient instantiates a k8s client and tests if it can list nodes
func GetK8sClient() (*kubernetes.Clientset, error) {
	k8sClient, err := loadClientFromServiceAccount()
	if err != nil {
		return nil, err
	} else if k8sClient == nil {
		return nil, ErrK8SApiAccountNotSet
	} else {
		_, err = k8sClient.CoreV1().Nodes().List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		return k8sClient, nil
	}
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
