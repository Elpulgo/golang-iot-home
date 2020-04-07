package main

import (
	"iot-home/endpoints"
	"iot-home/logger"
	"iot-home/netatmo"
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
	netatmoRest := new(netatmo.Rest)
	endpoints.Init(netatmo.New(netatmoRest))

	logger.Info("Listening on :3001 ...")
	error := http.ListenAndServe(":3001", nil)
	if error != nil {
		logger.Fatal(error.Error())
	}
}
