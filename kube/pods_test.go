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

	deployment := new(Deployment)
	err := c.CreateDeployment(deployment)

	assert.Nil(t, err)
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
