package kube

import "k8s.io/client-go/pkg/api/v1"

// GetServices retrieves services in the k8s cluster
func (c *Client) GetServices(namespace string, options v1.ListOptions) ([]v1.Service, error) {
	serviceList, err := c.k8sClient.Services(namespace).List(options)
	return serviceList.Items, err
}

// GetService retrieves specific service
func (c *Client) GetService(namespace string, name string) (*v1.Service, error) {
	svc, err := c.k8sClient.Services(namespace).Get(name)
	return svc, err
}

// CreateService will create a new service in the k8s cluster
// using the provided service spec
func (c *Client) CreateService(service *v1.Service) error {
	_, err := c.k8sClient.Services(service.Namespace).Create(service)

	return err
}

// UpdateService will update an existing service
func (c *Client) UpdateService(service *v1.Service) error {
	_, err := c.k8sClient.Services(service.Namespace).Update(service)

	return err
}
