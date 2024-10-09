// model/kubernetes_helpers.go
package model

import (
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ReadKubeConfig reads the kubeconfig file from the given path.
func ReadKubeConfig(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
	}
	return content, nil
}

// LoadKubeConfigFromBytes creates a Kubernetes clientset from kubeconfig content.
func LoadKubeConfigFromBytes(content []byte) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(content)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig content: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return clientset, nil
}
