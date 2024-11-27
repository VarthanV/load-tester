package models

import (
	"net/http"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Status string

const (
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

type Test struct {
	gorm.Model

	UUID                    uuid.UUID                       `gorm:"uniqueIndex" json:"uuid,omitempty"`
	URL                     string                          `json:"url,omitempty"`
	Method                  string                          `json:"method,omitempty"`
	Body                    datatypes.JSON                  `json:"body,omitempty"`
	Headers                 datatypes.JSONType[http.Header] `json:"headers,omitempty"`
	TargetUsers             int                             `json:"target_users,omitempty"`
	ReachPeakAfterInMinutes int                             `json:"reach_peak_after_in_minutes,omitempty"`
	UsersToStartWith        int                             `json:"users_to_start_with,omitempty"`
	TotalRequests           int32                           `json:"total_requests,omitempty"`
	SucceededRequests       int32                           `json:"succeeded_requests,omitempty"`
	FailedRequests          int32                           `json:"failed_requests,omitempty"`
	Report                  datatypes.JSON                  `json:"report,omitempty"`
}

func (t *Test) BeforeCreate(tx *gorm.DB) error {
	if t.UUID == uuid.Nil {
		t.UUID = uuid.New()
	}
	return nil
}
