package main

import (
	"fmt"
	"go-kafka-consumer/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

const (
	notificationsTopic = "notifications"
)

func main() {
	worker, err := ConnectConsumer([]string{"localhost:9092"})

	if err != nil {
		logger.Logger.Error(err.Error())
		panic(err)
	}

	consumer, err := worker.ConsumePartition(notificationsTopic, 0, sarama.OffsetOldest)

	if err != nil {
		logger.Logger.Error(err.Error())
		panic(err)
	}

	logger.Logger.Info("Consumer started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	messageCount := 0

	doneChan := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				logger.Logger.Error(err.Error())
			case msg := <-consumer.Messages():
				messageCount++
				logger.Logger.Info(fmt.Sprintf("Received message count: %d | Topic: %s | Message: %s", messageCount, string(msg.Topic), string(msg.Value)))
			case <-signalChan:
				logger.Logger.Info("Session ended")
				doneChan <- struct{}{}
			}
		}
	}()

	<-doneChan
	logger.Logger.Info(fmt.Sprintf("Processed %d messages", messageCount))

	err = worker.Close()

	if err != nil {
		logger.Logger.Error(err.Error())
		panic(err)
	}
}

func ConnectConsumer(brokerUrls []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerUrls, config)

	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	return conn, nil
}
