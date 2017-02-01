package kube

// APIInterface an interface for interacting with k8s cluster API
type APIInterface interface {
	Health() error
	GetPods() ([]Pod, error)
	CreateDeployment(*Deployment) error
	FindDeploymentByName(string) (*Deployment, error)
}
