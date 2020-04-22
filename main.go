package main

import (
	"fmt"
	"io"
	"iot-home/endpoints"
	"iot-home/hue"
	"iot-home/ioc"
	"iot-home/netatmo"
	"iot-home/wunderlist"
	"net/http"
	"os"

	"github.com/golobby/container"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

var (
	netatmoService    netatmo.Service
	wunderlistService wunderlist.Service
	port              string = ":3001"
)

func init() {
	initLogger()
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.WithField(".env", "Not found").Fatal("Will exit application ...")
	}

	ioc.Setup()
}

func main() {

	container.Make(&netatmoService)
	container.Make(&wunderlistService)
	endpoints.Init(netatmoService, wunderlistService)

	var hueRegistry hue.Registry

	container.Make(&hueRegistry)

	hueRegistry.Connect()

	logger.WithField("Port", port).Info("Started web server ...")
	error := http.ListenAndServe(port, nil)
	if error != nil {
		logger.WithError(error).Fatal()
	}
}
func initLogger() {
	file, err := os.OpenFile(
		"log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)

	if err != nil {
		fmt.Println(err)
	}

	logger.SetOutput(io.MultiWriter(file, os.Stdout))
}
