package portforwarder

import (
	"k8s.io/api/core/v1"
)

// isPodReady returns true if a pod is ready; false otherwise.
func isPodReady(pod *v1.Pod) bool {
	return isPodReadyConditionTrue(pod.Status)
}

// isPodReadyConditionTrue returns true if a pod is ready; false otherwise.
func isPodReadyConditionTrue(status v1.PodStatus) bool {
	condition := getPodReadyCondition(status)
	return condition != nil && condition.Status == v1.ConditionTrue
}

// getPodReadyCondition extracts the pod ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func getPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	_, condition := getPodCondition(&status, v1.PodReady)
	return condition
}

// getPodCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func getPodCondition(status *v1.PodStatus, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}
