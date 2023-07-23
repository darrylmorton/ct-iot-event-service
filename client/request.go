package client

import (
	"bytes"
	"fmt"
	"net/http"
)

type RequestOptions struct {
	Headers map[string]string
	Method  string
	Url     string
	Payload []byte
}

func GetHealthCheckRequest(requestOptions RequestOptions) *http.Response {
	request, newRequestErr := http.NewRequest(requestOptions.Method, requestOptions.Url, bytes.NewReader(requestOptions.Payload))

	if newRequestErr != nil {
		err := fmt.Errorf("GetHealthCheckRequest - newRequestErr %v", newRequestErr)
		fmt.Println(err.Error())
	}

	return App.GetHealthCheck(request)
}

func PutRequest(requestOptions RequestOptions) *http.Response {
	request, newRequestErr := http.NewRequest(requestOptions.Method, requestOptions.Url, bytes.NewReader(requestOptions.Payload))

	if newRequestErr != nil {
		err := fmt.Errorf("PutRequest - newRequestErr %v", newRequestErr)
		fmt.Println(err.Error())
	}

	return App.Put(request)
}

func GetAllRequest(requestOptions RequestOptions) *http.Response {
	request, newRequestErr := http.NewRequest(requestOptions.Method, requestOptions.Url, bytes.NewReader(requestOptions.Payload))

	if newRequestErr != nil {
		err := fmt.Errorf("GetAllRequest - newRequestErr %v", newRequestErr)
		fmt.Println(err.Error())
	}

	return App.GetAll(request)
}

func GetRequest(requestOptions RequestOptions) *http.Response {
	request, newRequestErr := http.NewRequest(requestOptions.Method, requestOptions.Url, bytes.NewReader(requestOptions.Payload))

	if newRequestErr != nil {
		err := fmt.Errorf("GetRequest - newRequestErr %v", newRequestErr)
		fmt.Println(err.Error())
	}

	return App.Get(request)
}
