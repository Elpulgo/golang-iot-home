package endpoints

import (
	"encoding/json"
	"fmt"
	"iot-home/netatmo"
	"net/http"
)

const STATIC_DIR = "./wwwroot/"
const API_PREFIX = "/api/v1/data/"

func Init(netatmoService netatmo.Service) {
	serveStaticContent()
	serveApiEndpoints(netatmoService)
}

func serveApiEndpoints(netatmoService netatmo.Service) {

	http.Header.Set(http.Header{}, "Content-Type", "application/json")

	http.Handle(API_PREFIX+"netatmo/current", serveNetatmoCurrent(netatmoService))
	// http.Handle(API_PREFIX+"netatmo/series", serveNetatmoSeries())
	// http.Handle(API_PREFIX+"wunderlist", serveWunderlist())
	// http.Handle(API_PREFIX+"hue", serveHue())
}

func serveNetatmoCurrent(service netatmo.Service) http.Handler {
	handlerFunc := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		current, error := service.GetCurrent()
		if error != nil {
			fmt.Println(error)
			responseWriter.Write([]byte("Failed to get current data from Netatmo!"))
			return
		}
		fmt.Println(current)

		json.NewEncoder(responseWriter).Encode(current)
	})

	return handlerFunc
}

// func serveNetatmoSeries() http.Handler {

// }

// func serveWunderlist() http.Handler {

// }

// func serveHue() http.Handler {

// }

func serveStaticContent() {
	http.Handle("/", index())
}

func index() http.Handler {
	return http.FileServer(http.Dir(STATIC_DIR))
}
