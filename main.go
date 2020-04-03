package main

import (
	"fmt"
	"iot-home/credentials"
	"iot-home/endpoints"
	"iot-home/logger"
	"iot-home/netatmo"
	"log"
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

	// token := new(netatmo.NetatmoOAuth)
	token, error := credentials.GetNetatmoOAuth()

	fmt.Println("Yippikayajdd")

	log.Println("Listening on :3001...")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
