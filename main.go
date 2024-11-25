package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/VarthanV/load-tester/pkg/tester"
)

func main() {
	// Define configuration options for the load tester
	driver, err := tester.New(
		// You can modify these options according to your use case
		tester.WithPeakConfig(50, 5*time.Minute, 10),             // Peak users and duration to reach peak
		tester.WithRequestConfig("http://example.com", nil, 200), // URL and body

	)
	if err != nil {
		log.Fatalf("Failed to create load tester: %v", err)
	}

	// Create a context with timeout for the test run (e.g., 20 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	// Start the load testing process
	fmt.Println("Starting the load testing...")
	go driver.Run(ctx)

	// Wait for the load test to finish
	select {
	case <-ctx.Done():
		fmt.Println("Load testing completed or timed out.")
	}

	// Print results

}
