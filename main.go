package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VarthanV/load-tester/pkg/tester"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	r.Run(":8060")

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
