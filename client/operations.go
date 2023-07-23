package client

import (
	"fmt"
	"net/http"
)

func (app *Application) Put(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("put - responseErr %+v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) GetAll(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("getAll - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}

func (app *Application) Get(req *http.Request) *http.Response {
	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		err := fmt.Errorf("get - responseErr %v", responseErr)
		fmt.Println(err.Error())
	}

	return response
}
