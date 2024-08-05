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
