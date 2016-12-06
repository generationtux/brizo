package kube

import (
	"github.com/generationtux/brizo/config"
	"github.com/urfave/cli"
)

// CLIFlags builds the config flags for kubernetes
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "k8s-config",
			Usage:       "path to kubeconfig file for access to external cluster (not needed if running brizo inside a k8s cluster)",
			Destination: &config.Kubernetes.ConfigFile,
			EnvVar:      "BRIZO_K8S_CONFIG",
		},
		cli.StringFlag{
			Name:        "k8s-context",
			Usage:       "context name to use when using an external kubeconfig file",
			Destination: &config.Kubernetes.Context,
			EnvVar:      "BRIZO_K8S_CONTEXT",
		},
	}
}
