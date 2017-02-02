package kube

import "k8s.io/client-go/pkg/api/v1"

// Pod spec
type Pod struct {
	Name string
}

// Deployment spec
type Deployment struct {
	Name string
	Pods []Pod
}

// PodOptions used to refine queries for pods
type PodOptions struct {
	Namespace string
	Name      string
	Labels    string
}

// GetPods retrieves pods in the k8s cluster
func (c *Client) GetPods(options PodOptions) ([]Pod, error) {
	k8sOptions := v1.ListOptions{
		LabelSelector: options.Labels,
	}

	podList, err := c.k8sClient.Pods(options.Namespace).List(k8sOptions)
	pods := make([]Pod, len(podList.Items))
	for i, pod := range podList.Items {
		pods[i] = Pod{
			Name: pod.Name,
		}
	}

	return pods, err
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
