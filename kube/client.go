package kube

import (
	"errors"
	"log"
	"os"

	"github.com/generationtux/brizo/config"
	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client client for interacting with Kubernetes cluster
type Client struct {
	external  bool
	k8sClient *kubernetes.Clientset
}

// New create new instance of kube client
func New() (*Client, error) {
	c := new(Client)
	c.external = config.Kubernetes.External
	err := c.initiateK8sClient()

	return c, err
}

// initializeK8sClient creates the k8s.io client
func (c *Client) initiateK8sClient() error {
	var k8sClient *kubernetes.Clientset
	var err error

	k8sClient, err = internalClient()
	if c.external {
		k8sClient, err = externalClient()
	}
	c.k8sClient = k8sClient

	return err
}

// Health checks the health of the kube client
func (c *Client) Health() error {
	if c.k8sClient == nil {
		return errors.New("Misconfigured k8s client.")
	}

	_, err := c.k8sClient.ServerVersion()
	if err != nil {
		return errors.New("Unable to reach k8s cluster.")
	}

	return nil
}

// internalClient builds the kube client configured to run inside a kube cluster
func internalClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		var client *kubernetes.Clientset
		return client, errors.New("Unable to connect to Kubernetes using in-cluster config.")
	}

	return kubernetes.NewForConfig(config)
}

// externalClient builds the kube client configured to run against an external kube cluster
func externalClient() (*kubernetes.Clientset, error) {
	var client *kubernetes.Clientset

	rules, err := getConfigLoadingRules()
	if err != nil {
		return client, err
	}

	overrides := getConfigOverrides()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		log.Println(err)
		return client, errors.New("Unable to connect to Kubernetes using external cluster config.")
	}

	return kubernetes.NewForConfig(config)
}

// getConfigLoadingRules returns k8s loading rules when external access is being used
func getConfigLoadingRules() (*clientcmd.ClientConfigLoadingRules, error) {
	rules := &clientcmd.ClientConfigLoadingRules{}
	path, err := getKubeConfigPath()
	rules.ExplicitPath = path

	if err != nil {
		return rules, err
	}

	return rules, nil
}

// getConfigOverrides returns k8s config when external access is being used
func getConfigOverrides() *clientcmd.ConfigOverrides {
	overrides := &clientcmd.ConfigOverrides{}

	if config.Kubernetes.Context != "" {
		overrides.CurrentContext = config.Kubernetes.Context
	}

	return overrides
}

// getKubeConfigPath will determine the path to the kubeconfig file
func getKubeConfigPath() (string, error) {
	file := config.Kubernetes.ConfigFile
	if file != "" {
		return validateFilePath(file)
	}

	file, err := validateFilePath("./kubeconfig")
	if err == nil {
		return file, err
	}

	return validateFilePath("~/.kube/config")
}

// validateFilePath will expand the provided path (looking for ~ home directory) and validate that the file exists
func validateFilePath(path string) (string, error) {
	actual, err := homedir.Expand(path)
	if err != nil {
		return actual, err
	}

	if _, err := os.Stat(actual); os.IsNotExist(err) {
		return actual, errors.New("kubeconfig file path does not exist")
	}

	return actual, nil
}
