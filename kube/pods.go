package kube

import (
	"fmt"
	"github.com/generationtux/brizo/resources"
	"k8s.io/client-go/pkg/api/v1"
)

// Pod is used to relate to Kubernetes pods
type Pod struct {
	Name string
}

// GetApplicationPods returns the pods running a provided application
func GetApplicationPods(app resources.Application) ([]Pod, error) {
	client, err := Client()
	if err != nil {
		return [0]Pod{}, error
	}

	kubePods := client.Pods("brizo")
	list := kubePods.List(v1.ListOptions{
		Labels: fmt.Sprintf("app_uuid=%v,brizoManaged=true", app.Name),
	})

	pods := []Pod{}
	for _, kubePod := range list.Items {
		pods = append(pods, Pod{Name: kubePod.getName()})
	}

	if len(pods) == 0 {
		pods = [0]Pod{}
	}

	return pods, nil
}
