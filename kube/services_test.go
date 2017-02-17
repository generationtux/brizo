package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/pkg/api/v1"
)

func TestGetServices(t *testing.T) {
	setup(t)
	c, _ := New()
	options := v1.ListOptions{
		LabelSelector: "appUUID=abc123",
	}
	services, err := c.GetServices("brizo", options)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
}
