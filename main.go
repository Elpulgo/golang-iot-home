package main

import (
	"fmt"
	"iot-home/credentials"
	"iot-home/endpoints"
	"iot-home/ioc"
	"iot-home/logger"
	"iot-home/netatmo"
	"net/http"

	"github.com/golobby/container"
	"github.com/joho/godotenv"
)

var (
	netatmoService netatmo.Service
	cred           credentials.CredentialsService
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.Fatal("No .env file found. Will exit application ...")
	}

	ioc.Setup()
}

func main() {

	container.Make(&netatmoService)
	container.Make(&cred)
	endpoints.Init(netatmoService)
	// netatmoRest := new(netatmoRest.Rest)
	// endpoints.Init(netatmo.New(netatmoRest))

	netatmoCred, err := cred.GetNetatmoOAuth()
	if err != nil {
		fmt.Println("wrong: " + err.Error())
	}

	fmt.Println(netatmoCred)

	logger.Info("Started web server, listening on :3001 ...")
	error := http.ListenAndServe(":3001", nil)
	if error != nil {
		logger.Fatal(error.Error())
	}
}
