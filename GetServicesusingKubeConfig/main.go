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

	// Define the namespace and service name.
	namespace := "default"
	serviceName := "opencost"

	// Fetch the service details from the Kubernetes cluster.
	fmt.Printf("Fetching service '%s' in namespace '%s'...\n", serviceName, namespace)
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Process and print the external IP or hostname if available.
	fmt.Println("Processing service ingress details...")
	for _, ingress := range service.Status.LoadBalancer.Ingress {
		if ingress.IP != "" {
			fmt.Printf("The external IP address for OpenCost is: %s\n", ingress.IP)
		} else if ingress.Hostname != "" {
			fmt.Printf("The external hostname for OpenCost is: %s\n", ingress.Hostname)
		} else {
			fmt.Println("No external IP address found for OpenCost service")
		}
	}

	fmt.Println("Program finished successfully.")
}
