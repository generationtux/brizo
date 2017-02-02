package kube

import (
	"errors"
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
)

// Pod spec
type Pod struct {
	Name string
}

// Deployment spec
type Deployment struct {
	Name        string
	Namespace   string
	Pods        []Pod
	MatchLabels map[string]string
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
func (c *Client) FindDeploymentByName(name, namespace string) (*Deployment, error) {
	k8sDeployment, err := c.k8sClient.Deployments(namespace).Get(name)
	if err != nil {
		return new(Deployment), err
	}

	deployment := new(Deployment)
	deployment.Name = k8sDeployment.Name
	deployment.Namespace = k8sDeployment.Namespace
	deployment.MatchLabels = k8sDeployment.Spec.Selector.MatchLabels

	pods, err := c.getPodsForDeployment(deployment)
	if err != nil {
		return new(Deployment), err
	}

	deployment.Pods = pods

	return deployment, nil
}

// getPodsForDeployment will retrieve the pods associated to the provided deployment
func (c *Client) getPodsForDeployment(deployment *Deployment) ([]Pod, error) {
	labelSelector := ""
	for key, value := range deployment.MatchLabels {
		labelSelector = fmt.Sprintf("%s,%s=%s", labelSelector, key, value)
	}
	labelSelector = labelSelector[1:] // remove first comma

	if labelSelector == "" {
		return []Pod{}, errors.New("No label selector specified for deployment")
	}

	return c.GetPods(PodOptions{
		Namespace: deployment.Namespace,
		Labels:    labelSelector,
	})
}
