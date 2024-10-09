// model/kubernetes.go
package model

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// FetchServices fetches all services from the Kubernetes cluster and returns OpenCost IP if found.
func FetchServices(clientset *kubernetes.Clientset) (string, error) {
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list services: %w", err)
	}

	for _, svc := range services.Items {
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				for _, port := range svc.Spec.Ports {
					if port.Port == 9090 {
						openCostIP := fmt.Sprintf("%s:%d", ingress.IP, port.Port)
						return openCostIP, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("OpenCost service with port 9090 not found")
}
