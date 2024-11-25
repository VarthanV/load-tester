package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/VarthanV/load-tester/pkg/tester"
)

func main() {
	driver, err := tester.New(
		tester.WithPeakConfig(20, 2*time.Minute, 10),
		tester.WithRequestConfig("http://example.com", nil, 200),
	)
	if err != nil {
		log.Fatalf("Failed to create load tester: %v", err)
	}

	fmt.Println("Starting the load testing...")
	driver.Run(context.Background())

	fmt.Println("Load testing completed or timed out.")

}
