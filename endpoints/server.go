package endpoints

import (
	"encoding/json"
	"iot-home/netatmo"
	"iot-home/wunderlist"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus"
)

const STATIC_DIR = "./wwwroot/"
const API_PREFIX = "/api/v1/data/"

func Init(netatmoService netatmo.Service, wunderlistService wunderlist.Service) {
	serveStaticContent()
	serveApiEndpoints(netatmoService, wunderlistService)
}

func serveApiEndpoints(netatmoService netatmo.Service, wunderlistService wunderlist.Service) {
	http.Header.Set(http.Header{}, "Content-Type", "application/json")

	http.Handle(API_PREFIX+"netatmo/current", serveNetatmoCurrent(netatmoService))
	http.Handle(API_PREFIX+"netatmo/series", serveNetatmoSeries(netatmoService))
	http.Handle(API_PREFIX+"wunderlist", serveWunderlist(wunderlistService))
	// http.Handle(API_PREFIX+"hue", serveHue())
}

func serveNetatmoCurrent(service netatmo.Service) http.Handler {
	handlerFunc := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		current, error := service.GetCurrent()
		if error != nil {
			logger.WithError(error).Error("Failed to get current data from Netatmo!")
			responseWriter.Write([]byte("Failed to get current data from Netatmo!"))
			return
		}
		json.NewEncoder(responseWriter).Encode(current)
	})

	return handlerFunc
}

func serveNetatmoSeries(service netatmo.Service) http.Handler {
	handlerFunc := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		current, error := service.GetHistory(time.Now().AddDate(0, 0, -3), time.Now().UTC())
		if error != nil {
			logger.WithError(error).Error("Failed to get historic data from Netatmo!")
			responseWriter.Write([]byte("Failed to get historic data from Netatmo!"))
			return
		}
		json.NewEncoder(responseWriter).Encode(current)
	})

	return handlerFunc
}

func serveWunderlist(service wunderlist.Service) http.Handler {
	handlerFunc := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		current, error := service.GetData()
		if error != nil {
			logger.WithError(error).Error("Failed to get tasks from Wunderlist!")
			responseWriter.Write([]byte("Failed to get tasks from Wunderlist!"))
			return
		}
		json.NewEncoder(responseWriter).Encode(current)
	})

	return handlerFunc
}

// func serveHue() http.Handler {

// }

func serveStaticContent() {
	http.Handle("/", index())
}

func index() http.Handler {
	return http.FileServer(http.Dir(STATIC_DIR))
}
