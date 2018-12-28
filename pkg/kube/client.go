package kube

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// GetServiceName retrieves the name of service in given namespace and labels
func GetServiceName(client corev1.ServicesGetter, namespace string, labels labels.Set) (string, error) {
	selector := labels.AsSelector()
	service, err := getFirstService(client, namespace, selector)
	if err != nil {
		return "", err
	}
	return service.ObjectMeta.GetName(), nil
}

func getFirstService(client corev1.ServicesGetter, namespace string, selector labels.Selector) (*v1.Service, error) {
	options := metav1.ListOptions{LabelSelector: selector.String()}
	services, err := client.Services(namespace).List(options)
	if err != nil {
		return nil, err
	}
	if len(services.Items) < 1 {
		return nil, fmt.Errorf("could not find service")
	}
	for _, s := range services.Items {
		return &s, nil
	}
	return nil, fmt.Errorf("could not find service")
}
