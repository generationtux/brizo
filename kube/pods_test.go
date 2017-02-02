package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPods(t *testing.T) {
	setup(t)
	c, _ := New()
	options := PodOptions{
		Namespace: "brizo",
		Labels:    "appUUID=abc123",
	}
	pods, err := c.GetPods(options)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(pods))
}

func TestCreateDeployment(t *testing.T) {
	setup(t)
	c, _ := New()

	deployment := &Deployment{
		Name:      "test-brizo-deploy-1",
		Namespace: "brizo",
		MatchLabels: map[string]string{
			"appUUID": "randomTestUUID1",
		},
	}
	err := c.CreateDeployment(deployment)
	assert.Nil(t, err)

	found, err := c.FindDeploymentByName(deployment.Name, deployment.Namespace)
	assert.Nil(t, err)
	assert.NotNil(t, found)
}

func TestFindDeploymentByName(t *testing.T) {
	setup(t)
	c, _ := New()
	deployment, err := c.FindDeploymentByName("product-mocker", "brizo")

	assert.Nil(t, err)
	assert.Equal(t, "product-mocker", deployment.Name)
	assert.Equal(t, "brizo", deployment.Namespace)
	assert.Equal(t, 3, len(deployment.Pods))
}
