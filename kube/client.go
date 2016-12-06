package kube

import (
	"github.com/generationtux/brizo/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client builds the kubernetes client
func Client() (*kubernetes.Clientset, error) {
	if isExternal() {
		return externalClient()
	}

	return internalClient()
}

// internalClient builds the kube client configured to run inside a kube cluster
func internalClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		client, _ := kubernetes.NewForConfig(config)
		return client, err
	}

	return kubernetes.NewForConfig(config)
}

// externalClient builds the kube client configured to run against an external kube cluster
func externalClient() (*kubernetes.Clientset, error) {
	rules := getConfigLoadingRules()
	overrides := getConfigOverrides()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()

	if err != nil {
		client, _ := kubernetes.NewForConfig(config)
		return client, err
	}

	return kubernetes.NewForConfig(config)
}

// isExternal determines if the k8s cluster is being accessed externally
func isExternal() bool {
	if config.Kubernetes.ConfigFile != "" {
		return true
	}

	return false
}

// getConfigLoadingRules returns k8s loading rules when external access is being used
func getConfigLoadingRules() *clientcmd.ClientConfigLoadingRules {
	rules := &clientcmd.ClientConfigLoadingRules{}

	if config.Kubernetes.ConfigFile != "" {
		rules.ExplicitPath = config.Kubernetes.ConfigFile
	}

	return rules
}

// getConfigOverrides returns k8s config when external access is being used
func getConfigOverrides() *clientcmd.ConfigOverrides {
	overrides := &clientcmd.ConfigOverrides{}

	if config.Kubernetes.Context != "" {
		overrides.CurrentContext = config.Kubernetes.Context
	}

	return overrides
}
