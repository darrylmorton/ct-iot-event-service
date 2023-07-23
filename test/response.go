package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PostEventResponse(res *http.Response) (int, Event) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err := fmt.Errorf("PostEventResponse - bodyCloseErr %v", err)
			fmt.Println(err.Error())
		}
	}(res.Body)

	statusCode := res.StatusCode

	body, responseBodyErr := io.ReadAll(res.Body)
	if responseBodyErr != nil {
		err := fmt.Errorf("PostEventResponse - responseBodyErr %v", responseBodyErr)
		fmt.Println(err.Error())
	}

	var responseObj Event
	responseObjUnmarshallErr := json.Unmarshal(body, &responseObj)
	if responseObjUnmarshallErr != nil {
		err := fmt.Errorf("PostEventResponse - responseObjUnmarshallErr %v", responseObjUnmarshallErr)
		fmt.Println(err.Error())
	}

	return statusCode, responseObj
}

func GetEventsResponse(res *http.Response) (int, []Event) {
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

	var responseObj []Event
	responseObjUnmarshallErr := json.Unmarshal(body, &responseObj)
	if responseObjUnmarshallErr != nil {
		err := fmt.Errorf("GetEventsResponse - responseObjUnmarshallErr %v", responseObjUnmarshallErr)
		fmt.Println(err.Error())
	}

	return statusCode, responseObj
}

func GetEventResponse(res *http.Response) (int, Event) {
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

	var responseObj Event
	responseObjUnmarshallErr := json.Unmarshal(body, &responseObj)
	if responseObjUnmarshallErr != nil {
		err := fmt.Errorf("GetEventResponse - responseObjUnmarshallErr %v", responseObjUnmarshallErr)
		fmt.Println(err.Error())
	}

	return statusCode, responseObj
}
