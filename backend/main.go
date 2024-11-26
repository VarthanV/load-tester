package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VarthanV/load-tester/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("load_tester.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("error in opening db ", err)
	}

	ctrl := controllers.Controller{DB: db}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	r.POST("/test", ctrl.ExecuteTest)

	r.Run(":8060")

	fmt.Println("Load testing completed or timed out.")

}
