package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	ID string `json:"id"` // Used to poll later and get report later
}

func (c *Controller) ExecuteTest(ctx *gin.Context) {
	var (
		request = CreateTestRequest{}
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Println("error in binding request ", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

}
