// main.go
package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"opencostProxy/controller"
	model "opencostProxy/models"
)

func main() {
	// Get the kubeconfig file path from the command line argument.
	kubeconfigPath := *flag.String("kubeconfig", "./kubeConfig/kubeconfig.yaml", "absolute path to the kubeconfig file")
	flag.Parse()

	// Load kubeconfig
	absKubeconfigPath, err := filepath.Abs(kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to get absolute path for kubeconfig: %s\n", err)
		return
	}

	kubeconfigContent, err := model.ReadKubeConfig(absKubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := model.LoadKubeConfigFromBytes(kubeconfigContent)
	if err != nil {
		panic(err.Error())
	}

	// Setup the controller
	openCostController := &controller.OpenCostController{Clientset: clientset}

	// Setup Gin router
	router := gin.Default()
	router.Use(model.CORS())

	// Routes
	router.GET("/model/allocation/compute", openCostController.FetchOpenCostService)

	// Start server
	router.Run(":8080")
}
