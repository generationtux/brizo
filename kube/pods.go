package kube

// Pod spec
type Pod struct {
	Name string
}

// Deployment spec
type Deployment struct {
	Name string
	Pods []Pod
}

// GetPods retrieves pods in the k8s cluster
func (c *Client) GetPods() ([]Pod, error) {
	return []Pod{}, nil
}

// CreateDeployment will create a new deployment in the k8s cluster
// using the provided deployment spec for configuration
func (c *Client) CreateDeployment(deployment *Deployment) error {
	return nil
}

// FindDeploymentByName will lookup a deployment in the k8s cluster using
// the provided name
func (c *Client) FindDeploymentByName(name string) (*Deployment, error) {
	deployment := new(Deployment)
	deployment.Name = name

	return deployment, nil
}
