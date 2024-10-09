// controller/opencost_controller.go
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"k8s.io/client-go/kubernetes"

	model "opencostProxy/models"
)

// OpenCostController struct to manage OpenCost interactions
type OpenCostController struct {
	Clientset *kubernetes.Clientset
}

// FetchOpenCostService fetches the OpenCost service IP and handles the request
func (ctrl *OpenCostController) FetchOpenCostService(c *gin.Context) {
	openCostIP, err := model.FetchServices(ctrl.Clientset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctrl.OpenCostEndPointHandler(c, openCostIP)
}

// OpenCostEndPointHandler handles the request to the OpenCost allocation compute endpoint
func (ctrl *OpenCostController) OpenCostEndPointHandler(c *gin.Context, openCostIP string) {
	url := fmt.Sprintf("http://%s/model/allocation/compute", openCostIP)

	// Extract query parameters
	queryParams := map[string]string{
		"window":      c.Query("window"),
		"aggregate":   c.Query("aggregate"),
		"includeIdle": c.Query("includeIdle"),
		"step":        c.Query("step"),
		"accumulate":  c.Query("accumulate"),
	}

	client := resty.New()
	resp, err := client.R().SetQueryParams(queryParams).Get(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, result)
}
