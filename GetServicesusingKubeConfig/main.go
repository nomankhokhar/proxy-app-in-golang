package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/go-resty/resty/v2"
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

// checkURL checks if the given URL is reachable.
func checkURL(url string) error {
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return fmt.Errorf("failed to GET URL %s: %w", url, err)
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		fmt.Printf("URL %s is reachable. Status Code: %d\n", url, resp.StatusCode())
		return nil
	}
	return fmt.Errorf("URL %s returned non-success status code: %d", url, resp.StatusCode())
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
	ingressName := "opencost-ingress"

	// Fetch the service details from the Kubernetes cluster.
	fmt.Printf("Fetching service '%s' in namespace '%s'...\n", serviceName, namespace)

	// Fetch the Ingress details from the Kubernetes cluster.
	fmt.Printf("Fetching ingress '%s' in namespace '%s'...\n", ingressName, namespace)
	ingress, err := clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}


	// Display the URL from the Ingress configuration.
	fmt.Println("Ingress URL: ", ingress.Spec.Rules[0].Host)
	url := fmt.Sprintf("http://%s/model/allocation/compute?window=7d&aggregate=namespace&includeIdle=true&step=1d&accumulate=false", ingress.Spec.Rules[0].Host)
	

	// Check if the URL is reachable.
	if err := checkURL(url); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("URL is working fine.")
	}

	fmt.Println("Program finished successfully.")
}
