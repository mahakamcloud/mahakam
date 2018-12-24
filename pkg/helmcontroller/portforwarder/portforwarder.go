package portforwarder

import (
	"fmt"

	"github.com/mahakamcloud/mahakam/pkg/kube"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

var (
	tillerPodLabels = labels.Set{"app": "helm", "name": "tiller"}
)

// New creates new and initialized tunnel
func New(namespace string, client kubernetes.Interface, config *rest.Config) (*kube.Tunnel, error) {
	podName, err := GetTillerPodName(client.CoreV1(), namespace)
	if err != nil {
		return nil, err
	}
	const tillerPort = 44134
	t := kube.NewTunnel(client.CoreV1().RESTClient(), config, namespace, podName, tillerPort)
	return t, t.ForwardPort()
}

// GetTillerPodName retrieves the name of tiller pod running in given namespace
func GetTillerPodName(client corev1.PodsGetter, namespace string) (string, error) {
	selector := tillerPodLabels.AsSelector()
	pod, err := getFirstRunningPod(client, namespace, selector)
	if err != nil {
		return "", err
	}
	return pod.ObjectMeta.GetName(), nil
}

func getFirstRunningPod(client corev1.PodsGetter, namespace string, selector labels.Selector) (*v1.Pod, error) {
	options := metav1.ListOptions{LabelSelector: selector.String()}
	pods, err := client.Pods(namespace).List(options)
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, fmt.Errorf("could not find tiller")
	}
	for _, p := range pods.Items {
		if isPodReady(&p) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("could not find ready tiller pod")
}
