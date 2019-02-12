package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func nodeWithStatus(nodeName string, nodePhase corev1.NodePhase) corev1.Node {
	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
		},
		Spec: corev1.NodeSpec{},
		Status: corev1.NodeStatus{
			Phase: nodePhase,
		},
	}
}

func TestValidateNodesWithNoNode(t *testing.T) {
	nodes := &corev1.NodeList{Items: []corev1.Node{}}

	_, err := validateNodes(nodes)
	assert.Error(t, err)
}

func TestValidateNodesWithNNode(t *testing.T) {
	var testcases = []struct {
		name  string
		nodes []corev1.Node

		expectFailures int
	}{
		{
			name: "Nodes start with not ready node",
			nodes: []corev1.Node{
				nodeWithStatus("test-node-1", corev1.NodePending),
				nodeWithStatus("test-node-2", corev1.NodeRunning),
			},
			expectFailures: 1,
		},
		{
			name: "Nodes end with not ready node",
			nodes: []corev1.Node{
				nodeWithStatus("test-node-1", corev1.NodeRunning),
				nodeWithStatus("test-node-2", corev1.NodePending),
			},
			expectFailures: 1,
		},
		{
			name: "Nodes are all not ready",
			nodes: []corev1.Node{
				nodeWithStatus("test-node-1", corev1.NodePending),
				nodeWithStatus("test-node-2", corev1.NodePending),
			},
			expectFailures: 2,
		},
		{
			name: "Nodes are all ready",
			nodes: []corev1.Node{
				nodeWithStatus("test-node-1", corev1.NodeRunning),
				nodeWithStatus("test-node-2", corev1.NodeRunning),
			},
			expectFailures: 0,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			nodes := &corev1.NodeList{
				Items: testcase.nodes,
			}

			failures, err := validateNodes(nodes)
			assert.NoError(t, err)
			assert.Equal(t, testcase.expectFailures, len(failures))
		})
	}
}
