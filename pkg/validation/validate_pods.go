package validation

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const namespace = "kube-system"

type ValidationError struct {
	name    string
	message string
}

func ValidatePods(kubeclient kubernetes.Interface) ([]*ValidationError, error) {
	pods, err := getPods(kubeclient.CoreV1(), namespace)
	if err != nil {
		return nil, err
	}
	verr, err := validatePods(pods, namespace)
	if err != nil {
		return nil, err
	}
	return verr, nil
}

func getPods(c corev1.PodsGetter, namespace string) (*v1.PodList, error) {
	options := metav1.ListOptions{}
	pods, err := c.Pods(namespace).List(options)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func validatePods(pods *v1.PodList, namespace string) ([]*ValidationError, error) {
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("pods in namespace %q not exist", namespace)
	}

	var failures []*ValidationError
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodSucceeded {
			continue
		}

		if pod.Status.Phase == v1.PodPending ||
			pod.Status.Phase == v1.PodFailed ||
			pod.Status.Phase == v1.PodUnknown {
			failures = append(failures, &ValidationError{
				name:    fmt.Sprintf("%q/%q", pod.Namespace, pod.Name),
				message: fmt.Sprintf("Pod %q in namespace %q is %s", pod.Name, pod.Namespace, pod.Status.Phase),
			})
			continue
		}

		var notready []string
		for _, container := range pod.Status.ContainerStatuses {
			if !container.Ready {
				notready = append(notready, container.Name)
			}
		}
		if len(notready) != 0 {
			failures = append(failures, &ValidationError{
				name:    fmt.Sprintf("%q/%q", pod.Namespace, pod.Name),
				message: fmt.Sprintf("Pod %q in namespace"),
			})
		}
	}

	return failures, nil
}
