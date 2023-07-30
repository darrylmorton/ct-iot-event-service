package test

import (
	"encoding/json"
	"fmt"
	"github.com/darrylmorton/ct-iot-event-service/internal/app"
	"io"
	"net/http"
)

func GetHealthCheckResponse(res *http.Response) (int, HealthCheck) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err := fmt.Errorf("GetHealthCheckResponse - bodyCloseErr %v", err)
			fmt.Println(err.Error())
		}
	}(res.Body)

	statusCode := res.StatusCode

	body, responseBodyErr := io.ReadAll(res.Body)
	if responseBodyErr != nil {
		err := fmt.Errorf("GetHealthCheckResponse - responseBodyErr %v", responseBodyErr)
		fmt.Println(err.Error())
	}

	var responseObj HealthCheck
	json.Unmarshal(body, &responseObj)

	return statusCode, responseObj
}

func PutEventResponse(res *http.Response) (int, app.Event) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err := fmt.Errorf("PutEventResponse - bodyCloseErr %v", err)
			fmt.Println(err.Error())
		}
	}(res.Body)

	statusCode := res.StatusCode

	body, responseBodyErr := io.ReadAll(res.Body)
	if responseBodyErr != nil {
		err := fmt.Errorf("PutEventResponse - responseBodyErr %v", responseBodyErr)
		fmt.Println(err.Error())
	}

	var responseObj app.Event
	json.Unmarshal(body, &responseObj)

	return statusCode, responseObj
}

func GetEventsResponse(res *http.Response) (int, []app.Event) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err := fmt.Errorf("GetEventsResponse - bodyCloseErr %v", err)
			fmt.Println(err.Error())
		}
	}(res.Body)

	statusCode := res.StatusCode

	body, responseBodyErr := io.ReadAll(res.Body)
	if responseBodyErr != nil {
		err := fmt.Errorf("GetEventsResponse - responseBodyErr %v", responseBodyErr)
		fmt.Println(err.Error())
	}

	var responseObj []app.Event
	json.Unmarshal(body, &responseObj)

	return statusCode, responseObj
}

func GetEventResponse(res *http.Response) (int, app.Event) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err := fmt.Errorf("GetEventResponse - bodyCloseErr %v", err)
			fmt.Println(err.Error())
		}
	}(res.Body)

	statusCode := res.StatusCode

	body, responseBodyErr := io.ReadAll(res.Body)
	if responseBodyErr != nil {
		err := fmt.Errorf("GetEventResponse - responseBodyErr %v", responseBodyErr)
		fmt.Println(err.Error())
	}

	var responseObj app.Event
	json.Unmarshal(body, &responseObj)

	return statusCode, responseObj
}
