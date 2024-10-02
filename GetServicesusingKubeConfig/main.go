package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var openCostIP string

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

// fetchServices fetches all services from the Kubernetes cluster and sets the openCostIP if found.
func fetchServices(clientset *kubernetes.Clientset) error {
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	for _, svc := range services.Items {
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				for _, port := range svc.Spec.Ports {
					if port.Port == 9090 {
						openCostIP = fmt.Sprintf("%s:%d", ingress.IP, port.Port)
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("OpenCost service with port 9090 not found")
}

func main() {
	// Get the kubeconfig file path from the command line argument.
	kubeconfigPath := *flag.String("kubeconfig", "./nomanProxy.yaml", "absolute path to the kubeconfig file")
	flag.Parse()

	// Read and load the kubeconfig.
	absKubeconfigPath, err := filepath.Abs(kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to get absolute path for kubeconfig: %s\n", err)
		return
	}

	kubeconfigContent, err := readKubeConfig(absKubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := loadKubeConfigFromBytes(kubeconfigContent)
	if err != nil {
		panic(err.Error())
	}

	// Fetch services and retrieve OpenCost IP
	if err := fetchServices(clientset); err != nil {
		fmt.Printf("Error fetching services: %s\n", err)
		return
	}

	fmt.Println("OpenCost service IP:", openCostIP)

	// Setup Gin router
	router := gin.Default()
	router.GET("/model/allocation/compute", func(c *gin.Context) {
		openCostEndPointHandler(c, openCostIP)
	})
	router.Run(":8080")
}

// openCostEndPointHandler handles the request to the OpenCost allocation compute endpoint.
func openCostEndPointHandler(c *gin.Context, openCostIP string) {
	
	url := fmt.Sprintf("http://%s/model/allocation/compute", openCostIP)
	
	// Extract query parameters
	window := c.Query("window")
	aggregate := c.Query("aggregate")
	includeIdle := c.Query("includeIdle")
	step := c.Query("step")
	accumulate := c.Query("accumulate")

	client := resty.New()


	resp, err := client.R().
		SetQueryParams(map[string]string{
			"window":      window,
			"aggregate":   aggregate,
			"includeIdle": includeIdle,
			"step":        step,
			"accumulate":  accumulate,
		}).Get(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err,
		})
		return
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to parse allocation summary response",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
