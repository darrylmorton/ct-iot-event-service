package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

func PostEvents(db *sql.DB, data []Event) (int, error) {
	if len(data) > 0 {
		stmt, err := db.Prepare(`
		INSERT INTO events (device_name, description, type, event, read)
		VALUES ($1, $2, $3, $4, $5)
	`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

		for _, event := range data {
			if _, err := stmt.Exec(event.DeviceName, event.Description, event.Type, event.Event, event.Read); err != nil {
				log.Fatal(err)
			}
		}

		return len(data), err
	} else {
		return 0, nil
	}
}

// GetMessages returns the parsed messages from SQS if any. If an error
// occurs that error will be returned.
func (q *Queue) GetMessages(db *sql.DB, waitTimeout int64) {
	params := sqs.ReceiveMessageInput{
		QueueUrl: aws.String(q.URL),
	}
	if waitTimeout > 0 {
		params.WaitTimeSeconds = aws.Int64(waitTimeout)
	}
	resp, err := q.Client.ReceiveMessage(&params)
	if err != nil {
		//return nil, fmt.Errorf("failed to get messages, %v", err)
		fmt.Errorf("failed to get messages, %v", err)
	}

	msgs := make([]Event, len(resp.Messages))

	for i, msg := range resp.Messages {
		parsedMsg := Event{}
		if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
			//return nil, fmt.Errorf("failed to unmarshal message, %v", err)
			fmt.Errorf("failed to unmarshal message, %v", err)
		}

		msgs[i] = parsedMsg
	}

	result, err := PostEvents(db, msgs)

	fmt.Println("*** *** GetMessages result", result)
	fmt.Println("*** *** GetMessages err", err)

	//return msgs, nil
}
