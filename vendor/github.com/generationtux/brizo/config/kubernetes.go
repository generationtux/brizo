package config

type kubeConfig struct {
	External   bool
	ConfigFile string
	Context    string
}

// Kubernetes holds configuration for the k8s cluster
var Kubernetes = kubeConfig{}
