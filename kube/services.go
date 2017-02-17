package kube

import "k8s.io/client-go/pkg/api/v1"

// GetServices retrieves services in the k8s cluster
func (c *Client) GetServices(namespace string, options v1.ListOptions) ([]v1.Service, error) {
	serviceList, err := c.k8sClient.Services(namespace).List(options)
	return serviceList.Items, err
}

// CreateService will create a new service in the k8s cluster
// using the provided service spec
func (c *Client) CreateService(service *v1.Service) error {
	_, err := c.k8sClient.Services(service.Namespace).Create(service)

	return err
}
