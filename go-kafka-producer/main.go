package main

import (
	"go-kafka-producer/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	dotEnvUtil utils.IDotEnvUtil = utils.DotEnvUtil
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
	routerGroup.GET("/test", TestAPI)

	router.Run(":" + port)
}

func TestAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
