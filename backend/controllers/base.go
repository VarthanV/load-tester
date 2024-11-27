package controllers

import (
	"github.com/VarthanV/load-tester/pkg/liveupdate"
	"gorm.io/gorm"
)

type Controller struct {
	DB      *gorm.DB
	Updates liveupdate.Updater
}
