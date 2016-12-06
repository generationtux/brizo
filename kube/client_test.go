package kube

import (
	"testing"

	"github.com/generationtux/brizo/config"
	"github.com/stretchr/testify/assert"
)

func TestIsExternal(t *testing.T) {
	assert.Equal(t, false, isExternal())

	config.Kubernetes.ConfigFile = "./foo"
	assert.Equal(t, true, isExternal())
}

func TestGetConfigLoadingRules(t *testing.T) {
	config.Kubernetes.ConfigFile = "./kubeconfig"
	rules := getConfigLoadingRules()
	assert.Equal(t, "./kubeconfig", rules.ExplicitPath)
}

func TestGetConfigOverrides(t *testing.T) {
	overrides := getConfigOverrides()
	assert.Equal(t, "", overrides.CurrentContext)

	config.Kubernetes.Context = "custom"
	overrides = getConfigOverrides()
	assert.Equal(t, "custom", overrides.CurrentContext)
}
