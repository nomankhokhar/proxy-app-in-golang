package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// readKubeConfig reads the kubeconfig file from the given path and returns its content as a byte slice.
func readKubeConfig(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
	}
	return content, nil
}

// loadKubeConfigFromBytes creates a Kubernetes clientset from kubeconfig content.
func loadKubeConfigFromBytes(content []byte) (*kubernetes.Clientset, error) {
	// Load the Kubernetes config from the kubeconfig content.
	config, err := clientcmd.RESTConfigFromKubeConfig(content)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig content: %w", err)
	}

	// Create a Kubernetes client from the config.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return clientset, nil
}

func main() {
	// Get the kubeconfig file path from the command line argument.
	kubeconfigPath := *flag.String("kubeconfig", "./nomanProxy.yaml", "absolute path to the kubeconfig file")
	flag.Parse()

	// Determine the absolute path to the kubeconfig file.
	absKubeconfigPath, err := filepath.Abs(kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to get absolute path for kubeconfig: %s\n", err)
		return
	}

	// Read the kubeconfig file content.
	kubeconfigContent, err := readKubeConfig(absKubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	// Load the Kubernetes clientset using the kubeconfig content.
	clientset, err := loadKubeConfigFromBytes(kubeconfigContent)
	if err != nil {
		panic(err.Error())
	}

	// Fetch all services in the cluster.
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Iterate through the services and print their external IPs and ports.
	fmt.Println("Listing External IPs and Ports of Services:")
	for _, svc := range services.Items {
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				for _, port := range svc.Spec.Ports {
					fmt.Printf("%s -> %s:%d\n", svc.Name, ingress.IP, port.Port)
				}
			}
		}
	}

	fmt.Println("Program finished successfully.")
}
