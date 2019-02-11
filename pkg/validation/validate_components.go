package validation

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func ValidateComponents(kubeclient kubernetes.Interface) ([]*ValidationError, error) {
	components, err := getComponents(kubeclient.CoreV1())
	if err != nil {
		return nil, err
	}
	return validateComponents(components)
}

func getComponents(c corev1.ComponentStatusesGetter) (*v1.ComponentStatusList, error) {
	options := metav1.ListOptions{}
	components, err := c.ComponentStatuses().List(options)
	if err != nil {
		return nil, err
	}
	return components, nil
}

func validateComponents(components *v1.ComponentStatusList) ([]*ValidationError, error) {
	if len(components.Items) == 0 {
		return nil, fmt.Errorf("kube-system components do not exist")
	}

	var failures []*ValidationError
	for _, component := range components.Items {
		for _, condition := range component.Conditions {
			if condition.Status != v1.ConditionTrue {
				failures = append(failures, &ValidationError{
					Name:    fmt.Sprintf("%q", component.Name),
					Message: fmt.Sprintf("Component %q is not healthy", component.Name),
				})
			}
		}
	}
	return failures, nil
}
