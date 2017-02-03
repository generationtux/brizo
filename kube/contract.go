package kube

import (
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// APIInterface an interface for interacting with k8s cluster API
type APIInterface interface {
	Health() error
	GetPods(string, v1.ListOptions) ([]v1.Pod, error)
	CreateDeployment(*v1beta1.Deployment) error
	DeleteDeployment(*v1beta1.Deployment) error
	FindDeploymentByName(string, string) (*v1beta1.Deployment, error)
}
