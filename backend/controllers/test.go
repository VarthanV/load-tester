package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/VarthanV/load-tester/models"
	"github.com/VarthanV/load-tester/pkg/liveupdate"
	"github.com/VarthanV/load-tester/pkg/tester"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

type GetUpdateResponse struct {
	Update *liveupdate.Update `json:"update"`
	Test   *models.Test       `json:"test"`
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
			c.Updates,
			tester.WithPeakConfig(
				request.TargetUsers,
				time.Duration(request.ReachPeakAferInMinutes*int(time.Minute)),
				request.UsersToStartWith),
			tester.WithRequestConfig(request.URL, nil, request.SuccessStatusCodes...),
			tester.WithDB(c.DB),
		)
		if err != nil {
			log.Printf("Failed to create load tester: %v\n", err)
		}

		log.Println("Starting for id ", t.UUID)
		driver.Run(ctx, t.UUID)

	}()

	ctx.JSON(http.StatusCreated, CreateTestResponse{
		ID: t.UUID,
	})
}

func (c *Controller) GetTest(ctx *gin.Context) {
	var (
		test = models.Test{}
	)
	testID := ctx.Param("id")
	if testID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid test id"))
		return
	}

	err := c.DB.Model(&models.Test{}).
		Where(&models.Test{
			UUID: uuid.MustParse(testID),
		}).Last(&test).Error
	if err != nil {
		log.Println("erorr in getting tests ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, test)
}

func (c *Controller) GetUpdate(ctx *gin.Context) {
	var (
		res  = GetUpdateResponse{}
		test *models.Test
	)
	testID := ctx.Param("id")
	if testID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid test id"))
		return
	}

	update, err := c.Updates.Get(uuid.MustParse(testID))
	if err != nil {
		logrus.Error("error in getting update ", err)

	}

	res.Update = update

	if update != nil &&
		update.TargetUsers ==
			update.TotalNumberofRequestsDone ||
		update == nil {
		// No need for report to be held from here on reached end
		// game
		c.Updates.Delete(uuid.MustParse(testID))
		// Report might be ready can fetch it
		err := c.DB.
			Model(&models.Test{}).
			Where(&models.Test{
				UUID: uuid.MustParse(testID),
			}).Last(&test).Error
		if err != nil {
			logrus.Error("unable to get test ")
		}
		res.Test = test
	}

	ctx.JSON(http.StatusOK, res)
}
