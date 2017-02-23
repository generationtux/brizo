package kube

import (
	"errors"
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// GetPods retrieves pods in the k8s cluster
func (c *Client) GetPods(namespace string, options v1.ListOptions) ([]v1.Pod, error) {
	podList, err := c.k8sClient.Pods(namespace).List(options)
	return podList.Items, err
}

// CreateDeployment will create a new deployment in the k8s cluster
// using the provided deployment spec for configuration
func (c *Client) CreateDeployment(deployment *v1beta1.Deployment) error {
	_, err := c.k8sClient.Deployments(deployment.Namespace).Create(deployment)

	return err
}

// DeleteDeployment will delete the provided deployment from k8s
func (c *Client) DeleteDeployment(deployment *v1beta1.Deployment) error {
	return c.k8sClient.Deployments(deployment.Namespace).Delete(deployment.Name, &v1.DeleteOptions{})
}

// CreateOrUpdateDeployment will gracefully apply updates to an existing
// deployment or create a fresh deployment
func (c *Client) CreateOrUpdateDeployment(deployment *v1beta1.Deployment) error {
	existingDeployment, err := c.FindDeploymentByName(deployment.Name, deployment.Namespace)

	if err != nil {
		return c.CreateDeployment(deployment)
	}

	existingDeployment.Spec = deployment.Spec

	_, err = c.k8sClient.Deployments(deployment.Namespace).Update(existingDeployment)

	return err
}

// FindDeploymentByName will lookup a deployment in the k8s cluster using
// the provided name
func (c *Client) FindDeploymentByName(name, namespace string) (*v1beta1.Deployment, error) {
	return c.k8sClient.Deployments(namespace).Get(name)
}

// GetPodsForDeployment will retrieve the pods associated to the provided deployment
func (c *Client) GetPodsForDeployment(deployment *v1beta1.Deployment) ([]v1.Pod, error) {
	labelSelector := ""
	for key, value := range deployment.Spec.Selector.MatchLabels {
		labelSelector = fmt.Sprintf("%s,%s=%s", labelSelector, key, value)
	}
	labelSelector = labelSelector[1:] // remove first comma

	if labelSelector == "" {
		return []v1.Pod{}, errors.New("No label selector specified for deployment")
	}

	return c.GetPods(deployment.Namespace, v1.ListOptions{
		LabelSelector: labelSelector,
	})
}
