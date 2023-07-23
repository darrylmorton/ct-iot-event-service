package client

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	App = Application{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
)

type EnvArgs struct {
	ClientProtocol string
	ClientHost     string
	ClientPort     string
	ClientUri      string
	ClientUrl      string
}

type Application struct {
	Client *http.Client
}

func init() {
	if os.Getenv("CLIENT_PROTOCOL") == "" {
		envErr := os.Setenv("CLIENT_PROTOCOL", "http")
		if envErr != nil {
			err := fmt.Errorf("CLIENT_PROTOCOL - envErr %+v", envErr)
			fmt.Println(err.Error())
		}
	}

	if os.Getenv("CLIENT_HOST") == "" {
		envErr := os.Setenv("CLIENT_HOST", "localhost")
		if envErr != nil {
			err := fmt.Errorf("CLIENT_HOST - envErr %+v", envErr)
			fmt.Println(err.Error())
		}
	}

	if os.Getenv("CLIENT_PORT") == "" {
		envErr := os.Setenv("CLIENT_PORT", "4000")
		if envErr != nil {
			err := fmt.Errorf("CLIENT_PORT - envErr %+v", envErr)
			fmt.Println(err.Error())
		}
	}
}

func CreateEnvs(clientProtocol string, clientHost string, clientPort string, clientUri string) EnvArgs {
	if clientProtocol == "" {
		clientProtocol = os.Getenv("CLIENT_PROTOCOL")
	}

	if clientHost == "" {
		clientHost = os.Getenv("CLIENT_HOST")
	}

	if clientPort == "" {
		clientPort = os.Getenv("CLIENT_PORT")
	}

	if clientUri == "" {
		clientUri = os.Getenv("CLIENT_URI")
	}

	clientUrl := fmt.Sprintf("%s://%s:%s/%s", clientProtocol, clientHost, clientPort, clientUri)

	return EnvArgs{
		ClientProtocol: clientProtocol,
		ClientHost:     clientHost,
		ClientPort:     clientPort,
		ClientUri:      clientUri,
		ClientUrl:      clientUrl,
	}
}
