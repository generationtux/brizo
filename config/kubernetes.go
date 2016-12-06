package config

type kubeConfig struct {
	ConfigFile string
	Context    string
}

// Kubernetes holds configuration for the k8s cluster
var Kubernetes = kubeConfig{}
