package main

import (
	"fmt"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/sdk"
)

// This is the script that will be compiled into a .so file.
// For example, to compile it:
//   go build -buildmode=plugin -o example_script.so example_script.go

// Process is the function the executor will look up in the .so plugin.
// It must match this exact signature if you plan to call it via reflection.
func Process(client *sdk.StreamFlexClient) error {
	// Example usage of the resources from the client
	fmt.Println("Received payload:", string(client.Payload))

	// Check DB connection (pseudocode)
	if client.DB != nil {
		// Use client.DB to query the database
		fmt.Println("DB connection is available")
	}

	// Check Redis client
	if client.RedisClient != nil {
		// For example: client.RedisClient.Set(...)
		fmt.Println("Redis connection is available")
	}

	// HTTP Client usage
	if client.HTTPClient != nil {
		fmt.Println("HTTP client is ready for requests")
	}

	// S3 / MinIO usage (via AWS session)
	if client.S3Client != nil {
		fmt.Println("S3 session is available")
		// e.g., s3.New(client.S3Client).PutObject(...)
	}

	// Kafka producer usage
	if client.KafkaProducer != nil {
		fmt.Println("Kafka Producer is available")
		// e.g., (*client.KafkaProducer)(topic, message)
	}

	// Your custom logic goes here...
	// e.g., parse JSON from client.Payload, call DB, call external API, etc.

	fmt.Println("example_script.go -> Process() completed successfully.")
	return nil
}
