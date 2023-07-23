package client

import (
	"fmt"
	"net/http"
)

func (app *Application) GetHealthCheck(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("getHealthCheck - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) PutEvent(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("putEvent - responseErr %+v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) GetEvents(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("getEvents - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) GetEvent(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("getEvent - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}
