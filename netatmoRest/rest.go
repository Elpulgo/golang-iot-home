package netatmoRest

import (
	"encoding/json"
	"io/ioutil"
	"iot-home/credentials"
	"iot-home/netatmo"
	"iot-home/utilities"
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"
)

type RestService interface {
	GetCurrent(result chan netatmo.CurrentResult)
}

type rest struct {
	service     RestService
	credentials credentials.CredentialsService
}

func New(credentials credentials.CredentialsService) RestService {
	return &rest{credentials: credentials}
}

func (rest *rest) GetCurrent(result chan netatmo.CurrentResult) {
	token, error := rest.credentials.GetNetatmoOAuth()

	if error != nil {
		logger.Error("Failed to get Netatmo current data from Netatmo API, access token not working")
		result <- netatmo.CurrentResult{Current: netatmo.Current{}, Error: error}
	}

	deviceId, deviceIdExists := os.LookupEnv("NETATMO_DEVICEID")

	if !deviceIdExists {
		logger.Error("Netatmo deviceid not set up in .env file! Can't fetch data from Netatmo API!")
		result <- netatmo.CurrentResult{Current: netatmo.Current{}, Error: error}
	}

	apiUrl := utilities.BuildStationUrl(token.AccessToken, deviceId).String()

	response, error := http.Get(apiUrl)

	if error != nil {
		logger.WithError(error).Error("Failed to get _current_ data from Netatmo API")
		result <- netatmo.CurrentResult{Current: netatmo.Current{}, Error: error}
	}

	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		logger.Error("Failed to parse body from Netatmo _current_ API")
		result <- netatmo.CurrentResult{Current: netatmo.Current{}, Error: error}
	}

	var currentData netatmo.Current

	json.Unmarshal(body, &currentData)

	result <- netatmo.CurrentResult{Current: currentData, Error: nil}
}
