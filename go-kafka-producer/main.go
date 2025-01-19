package main

import (
	"go-kafka-producer/controllers"
	"go-kafka-producer/utils"

	"github.com/gin-gonic/gin"
)

var (
	dotEnvUtil utils.IDotEnvUtil = utils.DotEnvUtil

	notificationController controllers.INotificationController = controllers.GetNotificationController()
)

func main() {
	port := dotEnvUtil.GetEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	routerGroup := router.Group("/api")
	routerGroup.POST("/notification", notificationController.CreateNewNotification)

	router.Run(":" + port)
}
