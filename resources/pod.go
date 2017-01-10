package resources

import (
	"fmt"

	"github.com/generationtux/brizo/kube"
	"k8s.io/client-go/pkg/api/v1"
)

// Pod is used to relate to Kubernetes pods
type Pod struct {
	Name string `json:"name,string"`
}

// PodRetrieval describes a method that can get a list of pods for a provided application UUID
type PodRetrieval func(string) ([]Pod, error)

// GetApplicationPods returns the pods running a provided application
func GetApplicationPods(UUID string) ([]Pod, error) {
	client, err := kube.Client()
	if err != nil {
		return []Pod{}, err
	}

	kubePods := client.Pods("brizo")
	list, err := kubePods.List(v1.ListOptions{
		LabelSelector: fmt.Sprintf("appUUID=%v, brizoManaged=true", UUID),
	})

	pods := []Pod{}
	for _, kubePod := range list.Items {
		pods = append(pods, Pod{Name: kubePod.GetName()})
	}

	if len(pods) == 0 {
		pods = make([]Pod, 0)
	}

	return pods, err
}
