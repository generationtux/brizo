package kube

import (
	"github.com/generationtux/brizo/resources"
)

// PodRetrieval describes a method that can get a list of pods for a provided application
type PodRetrieval func(resources.Application) []Pod
