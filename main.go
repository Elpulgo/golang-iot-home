package main

import (
	"iot-home/endpoints"
	"iot-home/logger"
	"iot-home/netatmo"
	"iot-home/netatmoRest"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.Fatal("No .env file found. Will exit application ...")
	}
}

func main() {
	netatmoRest := new(netatmoRest.Rest)
	endpoints.Init(netatmo.New(netatmoRest))

	logger.Info("Started web server, listening on :3001 ...")
	error := http.ListenAndServe(":3001", nil)
	if error != nil {
		logger.Fatal(error.Error())
	}
}
