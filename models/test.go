package models

import (
	"net/http"

	"github.com/VarthanV/load-tester/pkg/tester"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

	Report datatypes.JSONType[tester.Report] `json:"report,omitempty"`
}

type testRepo struct {
	db *gorm.DB
}
