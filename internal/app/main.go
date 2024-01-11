package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/darrylmorton/ct-iot-event-service/internal/data"
	"github.com/darrylmorton/ct-iot-event-service/internal/models"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const Version = "0.0.1"
const QueueName = "thing-payloads"

type EnvConfig struct {
	Port int
	Env  string
	Dsn  string
}

type Application struct {
	config EnvConfig
	logger *log.Logger
	models data.Models
}

func StartServer(serviceConfig *ServiceConfig) *http.Server {
	err := serviceConfig.DbClient.Ping()
	if err != nil {
		serviceConfig.Logger.Fatal(err)
	}
	serviceConfig.Logger.Printf("database connection pool established")

	app := &Application{
		config: serviceConfig.EnvConfig,
		logger: serviceConfig.Logger,
		models: data.NewModels(serviceConfig.DbClient),
	}

	go func() {
		result, err := GetMessages(context.TODO(), serviceConfig.SqsClient, serviceConfig.SqsReceiveMessageInput)
		if err != nil {
			serviceConfig.Logger.Printf("Error starting message consumer:%v\n", err)
		}

		messagesChannel := make(chan models.Event, 25)

		for _, item := range result.Messages {
			var unmarshalledMessage = models.Event{}
			var message = *item.Body

			err := json.Unmarshal([]byte(message), &unmarshalledMessage)
			if err != nil {
				serviceConfig.Logger.Printf("Error unmarshalling message:%v\n", err)
			} else {
				messagesChannel <- unmarshalledMessage
			}
		}

		_, err = serviceConfig.Models.Events.PostEvents(messagesChannel)
		if err != nil {
			serviceConfig.Logger.Printf("Error added message to database:%v\n", err)
		}

	}()

	addr := fmt.Sprintf(":%d", serviceConfig.EnvConfig.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	serviceConfig.Logger.Printf("starting %v server on %s\n", serviceConfig.EnvConfig, addr)

	return srv
}
