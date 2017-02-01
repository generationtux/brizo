package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPods(t *testing.T) {
	c, _ := New()
	pods, err := c.GetPods()

	assert.Nil(t, err)
	assert.Equal(t, []Pod{}, pods)
}

func TestCreateDeployment(t *testing.T) {
	c, _ := New()

	deployment := new(Deployment)
	err := c.CreateDeployment(deployment)

	assert.Nil(t, err)
}

func TestFindDeploymentByName(t *testing.T) {
	c, _ := New()
	deployment, err := c.FindDeploymentByName("foo")

	assert.Nil(t, err)
	assert.Equal(t, "foo", deployment.Name)
}
