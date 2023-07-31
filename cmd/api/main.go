package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/darrylmorton/ct-iot-event-service/internal/app"
)

func main() {
	sess := session.Must(session.NewSession())

	q := app.Queue{
		Client: sqs.New(sess),
		URL:    fmt.Sprintf("mockURL_%d", 0),
	}

	srv := app.StartServer(q)
	err := srv.ListenAndServe()
	fmt.Printf("The server has encountered an error: %v", err)
}
