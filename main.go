package main

import (
	"fmt"
	"io"
	"iot-home/credentials"
	"iot-home/endpoints"
	"iot-home/ioc"
	"iot-home/netatmo"
	"net/http"
	"os"

	"github.com/golobby/container"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

var (
	netatmoService netatmo.Service
	cred           credentials.CredentialsService
	port           string = ":3001"
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
	container.Make(&cred)
	endpoints.Init(netatmoService)

	netatmoCred, err := cred.GetNetatmoOAuth()
	if err != nil {
		logger.WithError(err)
	}

	logger.WithFields(logger.Fields{"netatmo cred": netatmoCred}).Info("hello")

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
