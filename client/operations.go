package client

import (
	"fmt"
	"net/http"
)

func (app *Application) Create(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("create - responseErr %+v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) GetAll(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("GetAll - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) Get(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("Get - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}
