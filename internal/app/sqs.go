package app

// GetMessages returns the parsed messages from SQS if any. If an error
// occurs that error will be returned.
//func (config *ServiceConfig) GetMessages(waitTimeout int64) {
//	params := sqs.ReceiveMessageInput{
//		QueueUrl: aws.String(config.SqsUrl),
//	}
//	if waitTimeout > 0 {
//		params.WaitTimeSeconds = int32(waitTimeout)
//	}
//	resp, err := config.SqsClient.ReceiveMessage(&params)
//	if err != nil {
//		//return nil, fmt.Errorf("failed to get messages, %v", err)
//		fmt.Errorf("failed to get messages, %v", err)
//	}
//
//	msgs := make([]models.Event, len(resp.Messages))
//
//	for i, msg := range resp.Messages {
//		parsedMsg := models.Event{}
//		if err := json.Unmarshal([]byte(*msg.Body), &parsedMsg); err != nil {
//			//return nil, fmt.Errorf("failed to unmarshal message, %v", err)
//			fmt.Errorf("failed to unmarshal message, %v", err)
//		}
//
//		msgs[i] = parsedMsg
//	}
//
//	result, err := config.Models.Events.PostEvents(msgs)
//
//	fmt.Println("*** *** GetMessages result", result)
//	fmt.Println("*** *** GetMessages err", err)
//
//	//return msgs, nil
//}
