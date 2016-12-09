package kube

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/generationtux/brizo/config"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigLoadingRules(t *testing.T) {
	ioutil.WriteFile("./kubeconfig", []byte("test"), 0644)
	config.Kubernetes.ConfigFile = "./kubeconfig"
	rules, err := getConfigLoadingRules()
	assert.Equal(t, "./kubeconfig", rules.ExplicitPath)
	assert.Equal(t, nil, err)
	os.Remove("./kubeconfig")

	config.Kubernetes.ConfigFile = "foobar"
	rules, err = getConfigLoadingRules()
	assert.NotNil(t, err)
}

func TestGetConfigOverrides(t *testing.T) {
	overrides := getConfigOverrides()
	assert.Equal(t, "", overrides.CurrentContext)

	config.Kubernetes.Context = "custom"
	overrides = getConfigOverrides()
	assert.Equal(t, "custom", overrides.CurrentContext)
}

func TestGetKubeConfigPath(t *testing.T) {
	// if config file is specified, validate exists
	config.Kubernetes.ConfigFile = "./client.go"
	path, err := getKubeConfigPath()
	assert.Equal(t, nil, err)
	assert.Equal(t, "./client.go", path)

	// if config file does not exist, try kubeconfig in current dir
	config.Kubernetes.ConfigFile = ""
	ioutil.WriteFile("./kubeconfig", []byte("test"), 0644)
	path, err = getKubeConfigPath()
	assert.Equal(t, nil, err)
	assert.Equal(t, "./kubeconfig", path)
	os.Remove("./kubeconfig")

	// if config file and ./kubeconfig does not exist, use ~/.kube/config
	path, err = getKubeConfigPath()
	expect, _ := homedir.Expand("~/.kube/config")
	assert.Equal(t, expect, path)
}

func TestValidateFilePath(t *testing.T) {
	// expands home dir
	path, _ := validateFilePath("~/foo")
	expect, _ := homedir.Expand("~/foo")
	assert.Equal(t, expect, path)

	// file not found
	_, err := validateFilePath("./foiaoijaeofij")
	assert.Equal(t, "kubeconfig file path does not exist", err.Error())
}
