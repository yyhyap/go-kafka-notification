package controllers

import (
	"encoding/json"
	"fmt"
	"go-kafka-producer/dto"
	"go-kafka-producer/logger"
	"net/http"
	"sync"

	"github.com/IBM/sarama"

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

	err = nc.PushToQueue(notificationsTopic, data)

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
	brokerUrls := []string{"localhost:9092"}
	producer, err := nc.ConnectProducer(brokerUrls)

	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)

	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset))

	return nil
}

func (nc *notificationControllerStruct) ConnectProducer(brokerUrls []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokerUrls, config)

	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	return conn, nil
}
