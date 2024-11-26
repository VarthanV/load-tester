package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VarthanV/load-tester/models"
	"github.com/VarthanV/load-tester/pkg/tester"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTestRequest struct {
	URL                    string            `json:"url"`
	Method                 string            `json:"method"`
	Body                   interface{}       `json:"body"`
	TargetUsers            int               `json:"target_users"`
	ReachPeakAferInMinutes int               `json:"reach_peak_afer_in_minutes"`
	Headers                map[string]string `json:"headers"`
	UsersToStartWith       int               `json:"users_to_start_with"`
	SuccessStatusCodes     []int             `json:"success_status_codes"`
}

type CreateTestResponse struct {
	ID uuid.UUID `json:"id"` // Used to poll later and get report later
}

func (c *Controller) ExecuteTest(ctx *gin.Context) {
	var (
		request = CreateTestRequest{}
		body    []byte
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Println("error in binding request ", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Body != nil {
		body, err = json.Marshal(request.Body)
		if err != nil {
			log.Println("error in marshalling body ", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	t := &models.Test{
		URL:                     request.URL,
		Method:                  request.Method,
		Body:                    body,
		UsersToStartWith:        request.UsersToStartWith,
		TargetUsers:             request.TargetUsers,
		ReachPeakAfterInMinutes: request.ReachPeakAferInMinutes,
		Status:                  models.StatusInProgress,
	}

	err = c.
		DB.
		Model(&models.Test{}).
		Create(t).Error
	if err != nil {
		log.Println("error in creating test ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {
		driver, err := tester.New(
			tester.WithPeakConfig(
				request.TargetUsers,
				time.Duration(request.ReachPeakAferInMinutes*int(time.Minute)),
				request.UsersToStartWith),
			tester.WithRequestConfig(request.URL, nil, 200),
		)
		if err != nil {
			log.Fatalf("Failed to create load tester: %v", err)
		}

		fmt.Println("Starting the load testing...")
		driver.Run(context.Background())

	}()

	ctx.JSON(http.StatusCreated, CreateTestResponse{
		ID: t.UUID,
	})
}
