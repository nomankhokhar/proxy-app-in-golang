package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-ClusterQueueLenght, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORS())

	router.GET("/ping", helloPing)
	router.GET("model/allocation/compute", paramAllocationHandler)

	router.Run(":8080")
}

func helloPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func paramAllocationHandler(c *gin.Context) {

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
		}).
		Get("http://206.189.246.89:9090/model/allocation/compute")

		
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

	c.JSON(http.StatusOK, result)
}
