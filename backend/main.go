package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VarthanV/load-tester/controllers"
	"github.com/VarthanV/load-tester/models"
	"github.com/VarthanV/load-tester/pkg/liveupdate"
	"github.com/gin-contrib/cors"
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

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache duration
	}

	ctrl := controllers.Controller{DB: db, Updates: liveupdate.New()}

	r.Use(cors.New(corsConfig))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	r.POST("/test", ctrl.ExecuteTest)

	r.GET("/test/:id", ctrl.GetTest)
	r.GET("/test/:id/updates", ctrl.GetUpdate)

	r.Run(":8060")

}
