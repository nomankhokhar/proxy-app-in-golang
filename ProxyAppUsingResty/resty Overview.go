package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// Middleware function to log request details
func logRequest(c *resty.Client, r *resty.Request) error {
	log.Printf("Request URL: %s", r.URL)
	return nil
}

// Middleware function to log response details
func logResponse(c *resty.Client, r *resty.Response) error {
	log.Printf("Response Status Code: %d", r.StatusCode())
	return nil
}

func main() {
	// Create a new Resty client
	client := resty.New()

	// Configure the client
	client.SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	// Add middlewares
	client.OnBeforeRequest(logRequest)
	client.OnAfterResponse(logResponse)

	// Create (POST) a new resource
	createResponse, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{"title": "foo", "body": "bar", "userId": 1}).
		Post("https://jsonplaceholder.typicode.com/posts")

	if err != nil {
		log.Fatalf("Create (POST) Error: %v", err)
	}
	fmt.Println("Create (POST) Response:", createResponse)

	// Read (GET) a resource
	readResponse, err := client.R().
		SetHeader("Accept", "application/json").
		Get("https://jsonplaceholder.typicode.com/posts/1")

	if err != nil {
		log.Fatalf("Read (GET) Error: %v", err)
	}
	fmt.Println("Read (GET) Response:", readResponse)

	// Update (PUT) a resource
	updateResponse, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{"id": 1, "title": "foo", "body": "bar", "userId": 1}).
		Put("https://jsonplaceholder.typicode.com/posts/1")

	if err != nil {
		log.Fatalf("Update (PUT) Error: %v", err)
	}
	fmt.Println("Update (PUT) Response:", updateResponse)

	// Delete (DELETE) a resource
	deleteResponse, err := client.R().
		Delete("https://jsonplaceholder.typicode.com/posts/1")

	if err != nil {
		log.Fatalf("Delete (DELETE) Error: %v", err)
	}
	fmt.Println("Delete (DELETE) Response:", deleteResponse)
}
