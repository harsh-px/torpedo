package k8sutils

import (
	"bytes"
	"text/template"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	/*
	"k8s.io/client-go/pkg/api"
	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
	"fmt"
	"github.com/Sirupsen/logrus"
	*/
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)
const k8sMasterLabelKey="node-role.kubernetes.io/master"

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

// DeployApp deploys given spec
func DeployApp(filepath string, params interface{}) error {
/*	specText, err := parseYAMLTemplate(filepath, params)
	if err != nil {
		return &ErrFailedToParseYAML{
			Path: filepath,
			Cause: err.Error(),
		}
	}

	decode := api.Codecs.UniversalDeserializer().Decode
	specObj, _, err := decode([]byte(specText), nil, nil)
	if err != nil {
		return &ErrFailedToApplySpec{
			Path:filepath,
			Cause: err.Error(),
		}
	}

	switch specObj.GetObjectKind().GroupVersionKind().Kind {
	case "Deployment":
		logrus.Infof("Found a deployment!")
		err = doDeployment(specObj.(*v1beta1.Deployment))
	case "PersistentVolumeClaim":
		logrus.Infof("Found a PersistentVolumeClaim!")
	case "StorageClass":
		logrus.Infof("Found a PersistentVolumeClaim!")
	default:
		return &ErrFailedToApplySpec{
			Path:filepath,
			Cause: fmt.Sprintf("Found unhandled type: %#v", specObj),
		}
	}


	fmt.Printf("%#v\n", specObj)*/
	return nil
}

// IsNodeMaster returns true if given node is a kubernetes master node
func IsNodeMaster(node v1.Node) bool {
	_, ok := node.Labels[k8sMasterLabelKey]
	return ok
}


func doDeployment(deployment *v1beta1.Deployment) error {
	return nil
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

// parseYAMLTemplate parses given file in GTPL format using the given params
func parseYAMLTemplate(filepath string, params interface{}) (string, error) {
	t, err := template.ParseFiles(filepath)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	if err = t.Execute(&result, params); err != nil {
		return "", err
	}

	return result.String(), nil
}