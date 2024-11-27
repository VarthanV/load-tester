package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VarthanV/load-tester/controllers"
	"github.com/VarthanV/load-tester/models"
	"github.com/VarthanV/load-tester/pkg/liveupdate"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	r := gin.Default()

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(sqlite.Open("load_tester.db"), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Fatal("error in opening db ", err)
	}

	err = db.AutoMigrate(&[]models.Test{})
	if err != nil {
		log.Fatal("unable to migrate tables ", err)
	}

	// Intiailize live updater

	ctrl := controllers.Controller{DB: db, Updates: liveupdate.New()}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	r.POST("/test", ctrl.ExecuteTest)

	r.GET("/test/:id", ctrl.GetTest)
	r.GET("/test/:id/updates", ctrl.GetUpdate)

	r.Run(":8060")

}
