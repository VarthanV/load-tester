package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/VarthanV/load-tester/config"
	"github.com/VarthanV/load-tester/controllers"
	"github.com/VarthanV/load-tester/models"
	"github.com/VarthanV/load-tester/pkg/liveupdate"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	r := gin.Default()

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatal("error in loading cfg ", err)
	}

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(sqlite.Open(cfg.Database.DatabaseName), &gorm.Config{
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
		AllowOrigins: strings.Split(cfg.Server.AllowedHosts, ","),
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{
			"Content-Length",
		},
		MaxAge: 12 * time.Hour, // Cache duration
	}

	ctrl := controllers.Controller{DB: db, Updates: liveupdate.New()}

	r.Use(cors.New(corsConfig))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	testsGroup := r.Group("/tests")

	testsGroup.POST("", ctrl.ExecuteTest)

	testsGroup.GET("/:id", ctrl.GetTest)
	testsGroup.GET("/:id/updates", ctrl.GetUpdate)

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))

}
