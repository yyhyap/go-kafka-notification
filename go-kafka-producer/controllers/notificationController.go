package controllers

import (
	"encoding/json"
	"go-kafka-producer/dto"
	"go-kafka-producer/logger"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	notificationsTopic = "notifications"
)

var (
	notificationController *notificationControllerStruct

	once sync.Once
)

type INotificationController interface {
	CreateNewNotification(c *gin.Context)
	PushToQueue(topic string, message []byte) error
}

type notificationControllerStruct struct{}

func GetNotificationController() INotificationController {
	if notificationController == nil {
		once.Do(func() {
			notificationController = &notificationControllerStruct{}
		})
	}

	return notificationController
}

func (nc *notificationControllerStruct) CreateNewNotification(c *gin.Context) {
	var newNotification dto.Notification

	err := c.BindJSON(&newNotification)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, dto.HttpResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// push comment to queue
	data, err := json.Marshal(newNotification)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, dto.HttpResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	nc.PushToQueue(notificationsTopic, data)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, dto.HttpResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.HttpResponse{
		Success: true,
		Message: "Notification created successfully",
		Data:    newNotification,
	})
}

func (nc *notificationControllerStruct) PushToQueue(topic string, message []byte) error {
	return nil
}
