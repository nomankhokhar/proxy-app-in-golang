package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func main() {
	router := gin.Default()

	router.GET("/ping", pingHandler)
	router.GET("/allocation", allocationHandler)
	router.GET("/allocation/:param1/:param2", paramAllocationHandler)

	router.Run(":8080")
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func allocationHandler(c *gin.Context) {
	client := resty.New()

	// Make a request to the other API
	resp, err := client.R().Get("http://localhost:3000/allocation/summary")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to fetch allocation summary",
		})
		return
	}

	var result map[string]interface{}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse allocation summary response",
		})
		return
	}

	c.JSON(resp.StatusCode(), result)
}

func paramAllocationHandler(c *gin.Context) {
	param1 := c.Param("param1")
	param2 := c.Param("param2")
	param3 := c.Param("param3")

	client := resty.New()

	resp, err := client.R().
		SetQueryParam("param1", param1).
		SetQueryParam("param2", param2).
		SetQueryParam("param3", param3).
		Get("http://localhost:3000/allocation/summary")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to fetch allocation summary",
		})

		return
	}

	var result map[string]interface{}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to perase allocation summary response",
		})
		return
	}

	c.JSON(resp.StatusCode(), result)

}
