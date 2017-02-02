package kube

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/generationtux/brizo/config"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	var status int
	if os.Getenv("BRIZO_TEST_KUBE") == "true" {
		setupTests(os.Getenv("BRIZO_TEST_KUBE_CONTEXT"))
		status = m.Run()
	} else {
		status = 0
		fmt.Println("Skipping kube tests. Set BRIZO_TEST_KUBE=true to run tests.")
	}
	os.Exit(status)
}

func setupTests(kubeContext string) {
	if kubeContext == "" {
		kubeContext = "minikube"
	}

	config.Kubernetes.External = true
	config.Kubernetes.Context = kubeContext
}

func TestNewClient(t *testing.T) {
	c, _ := New()
	assert.True(t, c.external)
	assert.NotEmpty(t, c.k8sClient)
}

// Expects the connected k8s cluster to be external
func TestClientHealth(t *testing.T) {
	c, _ := New()
	assert.Nil(t, c.Health())

	config.Kubernetes.External = false
	c, _ = New()
	assert.NotNil(t, c.Health())
}

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
	assert.Nil(t, err)
	assert.Equal(t, "./kubeconfig", path)
	os.Remove("./kubeconfig")

	// if config file and ./kubeconfig does not exist, use ~/.kube/config
	path, _ = getKubeConfigPath()
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
