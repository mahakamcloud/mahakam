package validation

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func ValidateNodes(kubeclient kubernetes.Interface) ([]*ValidationError, error) {
	nodes, err := getNodes(kubeclient.CoreV1())
	if err != nil {
		return nil, err
	}
	return validateNodes(nodes)
}

func getNodes(c corev1.NodesGetter) (*v1.NodeList, error) {
	options := metav1.ListOptions{}
	nodes, err := c.Nodes().List(options)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func validateNodes(nodes *v1.NodeList) ([]*ValidationError, error) {
	if len(nodes.Items) == 0 {
		return nil, fmt.Errorf("kubernetes nodes do not exist")
	}

	var failures []*ValidationError
	for _, node := range nodes.Items {
		if node.Status.Phase == v1.NodeTerminated {
			continue
		}

		if node.Status.Phase == v1.NodePending {
			failures = append(failures, &ValidationError{
				Name:    fmt.Sprintf("%q", node.Name),
				Message: fmt.Sprintf("Node %q is %s", node.Name, node.Status.Phase),
			})
			continue
		}
	}

	return failures, nil
}
