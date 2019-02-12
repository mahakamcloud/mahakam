package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func componentConditionWithStatus(status corev1.ConditionStatus) corev1.ComponentCondition {
	return corev1.ComponentCondition{
		Status: status,
	}
}

func componentStatusListWithCondition(componentConditions []corev1.ComponentCondition) *corev1.ComponentStatusList {
	return &corev1.ComponentStatusList{
		Items: []corev1.ComponentStatus{
			{
				Conditions: componentConditions,
			},
		},
	}
}

func TestValidateComponentsWithNoComponent(t *testing.T) {
	components := &corev1.ComponentStatusList{Items: []corev1.ComponentStatus{}}

	_, err := validateComponents(components)
	assert.Error(t, err)
}

func TestValidateComponentsWithOneComponent(t *testing.T) {
	var testcases = []struct {
		name            string
		conditionStatus corev1.ConditionStatus
		expectFailures  int
	}{
		{
			name:            "Components include unknown status",
			conditionStatus: corev1.ConditionUnknown,
			expectFailures:  1,
		},
		{
			name:            "Components include not ready status",
			conditionStatus: corev1.ConditionFalse,
			expectFailures:  1,
		},
		{
			name:            "Components are all ready",
			conditionStatus: corev1.ConditionTrue,
			expectFailures:  0,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			components := componentStatusListWithCondition(
				[]corev1.ComponentCondition{
					componentConditionWithStatus(testcase.conditionStatus),
				},
			)

			failures, err := validateComponents(components)
			assert.NoError(t, err)
			assert.Equal(t, testcase.expectFailures, len(failures))
		})
	}
}
