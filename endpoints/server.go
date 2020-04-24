package endpoints

import (
	"encoding/json"
	"iot-home/hue"
	"iot-home/netatmo"
	"iot-home/wunderlist"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus"
)

const STATIC_DIR = "./wwwroot/"
const API_PREFIX = "/api/v1/data/"

func Init(
	netatmoService netatmo.Service,
	wunderlistService wunderlist.Service,
	hueService hue.Service) {
	serveStaticContent()
	serveApiEndpoints(netatmoService, wunderlistService, hueService)
}

func serveApiEndpoints(
	netatmoService netatmo.Service,
	wunderlistService wunderlist.Service,
	hueService hue.Service) {
	http.Header.Set(http.Header{}, "Content-Type", "application/json")

	http.Handle(API_PREFIX+"netatmo/current", serveNetatmoCurrent(netatmoService))
	http.Handle(API_PREFIX+"netatmo/series", serveNetatmoSeries(netatmoService))
	http.Handle(API_PREFIX+"wunderlist", serveWunderlist(wunderlistService))
	http.Handle(API_PREFIX+"hue", serveHue(hueService))
}

func serveNetatmoCurrent(service netatmo.Service) http.Handler {
	handlerFunc := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		current, error := service.GetCurrent()
		if error != nil {
			handleError(error, "Failed to get current data from Netatmo!", responseWriter)
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
			handleError(error, "Failed to get historic data from Netatmo!", responseWriter)
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
			handleError(error, "Failed to get tasks from Wunderlist!", responseWriter)
			return
		}
		json.NewEncoder(responseWriter).Encode(current)
	})

	return handlerFunc
}

func serveHue(service hue.Service) http.Handler {
	handlerFunc := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		hueLightDto, error := service.GetLights()
		if error != nil {
			handleError(error, "Failed to get lights from Philiphs Hue API!", responseWriter)
			return
		}
		json.NewEncoder(responseWriter).Encode(hueLightDto)
	})

	return handlerFunc
}
func handleError(error error, message string, responseWriter http.ResponseWriter) {
	logger.WithError(error).Error(message)
	responseWriter.WriteHeader(500)
	responseWriter.Write([]byte(message))
}

func serveStaticContent() {
	http.Handle("/", index())
}

func index() http.Handler {
	return http.FileServer(http.Dir(STATIC_DIR))
}
