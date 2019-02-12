package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func podWithStatus(podName, namespace string, podPhase corev1.PodPhase, containerReadyStatus bool) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{},
		Status: corev1.PodStatus{
			Phase: podPhase,
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Ready: containerReadyStatus,
				},
			},
		},
	}
}

func TestValidatePodsWithNoPod(t *testing.T) {
	pods := &corev1.PodList{Items: []corev1.Pod{}}

	_, err := validatePods(pods, "test-namespace")
	assert.Error(t, err)
}

func TestValidatePodsWithOnePod(t *testing.T) {
	var testcases = []struct {
		name                 string
		podPhase             corev1.PodPhase
		containerReadyStatus bool
		expectFailures       int
	}{
		{
			name:                 "Pods include terminating pod",
			podPhase:             corev1.PodSucceeded,
			containerReadyStatus: false,
			expectFailures:       0,
		},
		{
			name:                 "Pods include pending pod",
			podPhase:             corev1.PodPending,
			containerReadyStatus: false,
			expectFailures:       1,
		},
		{
			name:                 "Pods include failed pod",
			podPhase:             corev1.PodFailed,
			containerReadyStatus: false,
			expectFailures:       1,
		},
		{
			name:                 "Pods include unknown pod",
			podPhase:             corev1.PodUnknown,
			containerReadyStatus: false,
			expectFailures:       1,
		},
		{
			name:                 "Pods include pod with non-ready container",
			podPhase:             corev1.PodRunning,
			containerReadyStatus: false,
			expectFailures:       1,
		},
		{
			name:                 "Pods are all ready",
			podPhase:             corev1.PodRunning,
			containerReadyStatus: true,
			expectFailures:       0,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			pods := &corev1.PodList{
				Items: []corev1.Pod{
					podWithStatus("test-pod", "test-namespace", testcase.podPhase, testcase.containerReadyStatus),
				},
			}

			failures, err := validatePods(pods, "test-namespace")
			assert.NoError(t, err)
			assert.Equal(t, testcase.expectFailures, len(failures))
		})
	}
}

func TestValidatePodsWithNPods(t *testing.T) {
	var testcases = []struct {
		name string
		pods []corev1.Pod

		expectFailures int
	}{
		{
			name: "Pods start with failed pod",
			pods: []corev1.Pod{
				podWithStatus("test-pod-1", "test-namespace", corev1.PodFailed, false),
				podWithStatus("test-pod-2", "test-namespace", corev1.PodRunning, true),
			},
			expectFailures: 1,
		},
		{
			name: "Pods end with failed pod",
			pods: []corev1.Pod{
				podWithStatus("test-pod-1", "test-namespace", corev1.PodRunning, true),
				podWithStatus("test-pod-2", "test-namespace", corev1.PodFailed, false),
			},
			expectFailures: 1,
		},
		{
			name: "Pods include pod with non-ready container",
			pods: []corev1.Pod{
				podWithStatus("test-pod-1", "test-namespace", corev1.PodRunning, false),
				podWithStatus("test-pod-2", "test-namespace", corev1.PodRunning, true),
			},
			expectFailures: 1,
		},
		{
			name: "Pods are all failing",
			pods: []corev1.Pod{
				podWithStatus("test-pod-1", "test-namespace", corev1.PodFailed, false),
				podWithStatus("test-pod-2", "test-namespace", corev1.PodFailed, false),
			},
			expectFailures: 2,
		},
		{
			name: "Pods are all ready",
			pods: []corev1.Pod{
				podWithStatus("test-pod-1", "test-namespace", corev1.PodRunning, true),
				podWithStatus("test-pod-2", "test-namespace", corev1.PodRunning, true),
			},
			expectFailures: 0,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			pods := &corev1.PodList{
				Items: testcase.pods,
			}

			failures, err := validatePods(pods, "test-namespace")
			assert.NoError(t, err)
			assert.Equal(t, testcase.expectFailures, len(failures))
		})
	}
}
