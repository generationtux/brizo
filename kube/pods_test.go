package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	metav1 "k8s.io/client-go/pkg/apis/meta/v1"
)

func TestGetPods(t *testing.T) {
	setup(t)
	c, _ := New()
	options := v1.ListOptions{
		LabelSelector: "appUUID=abc123",
	}
	pods, err := c.GetPods("brizo", options)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(pods))
}

func TestCreateAndDeleteDeployment(t *testing.T) {
	setup(t)
	c, _ := New()

	container := v1.Container{}
	container.Name = "app"
	container.Image = "nexmill/mock-go-api:latest"
	container.ImagePullPolicy = "Always"
	container.Ports = []v1.ContainerPort{
		v1.ContainerPort{
			Protocol:      v1.ProtocolTCP,
			ContainerPort: 8000,
		},
	}

	podTemplate := v1.PodTemplateSpec{}
	podTemplate.Labels = map[string]string{
		"appUUID":      "randomTestUUID1",
		"brizoManaged": "true",
	}
	podTemplate.Spec.Containers = []v1.Container{container}

	replicas := int32(2)
	spec := v1beta1.DeploymentSpec{}
	spec.Replicas = &replicas
	spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"appUUID": "randomTestUUID1",
		},
	}
	spec.Template = podTemplate

	deployment := new(v1beta1.Deployment)
	deployment.Name = "test-brizo-deploy-1"
	deployment.Namespace = "brizo"
	deployment.Spec = spec

	err := c.CreateDeployment(deployment)
	assert.Nil(t, err)

	found, err := c.FindDeploymentByName(deployment.Name, deployment.Namespace)
	assert.Nil(t, err)
	assert.NotNil(t, found)

	// cleanup
	c.DeleteDeployment(deployment)
}

func TestFindDeploymentByName(t *testing.T) {
	setup(t)
	c, _ := New()
	deployment, err := c.FindDeploymentByName("product-mocker", "brizo")
	assert.Nil(t, err)
	assert.Equal(t, "product-mocker", deployment.Name)
	assert.Equal(t, "brizo", deployment.Namespace)

	deploymentPods, err := c.GetPodsForDeployment(deployment)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(deploymentPods))
}
