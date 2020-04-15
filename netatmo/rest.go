package netatmo

import (
	"encoding/json"
	"io/ioutil"
	"iot-home/credentials"
	"iot-home/models"
	"iot-home/utilities"
	"net/http"
	"os"
	"time"

	logger "github.com/sirupsen/logrus"
)

type RestService interface {
	// TODO: Return our datamodel instead...
	GetCurrent(result chan CurrentResult)
	GetHistory(start time.Time, end time.Time, result chan HistoricResult)
}

type rest struct {
	service     RestService
	credentials credentials.CredentialsService
}

func NewRestService(credentials credentials.CredentialsService) RestService {
	return &rest{credentials: credentials}
}

func (rest *rest) GetHistory(start time.Time, end time.Time, result chan HistoricResult) {
	token, error := rest.credentials.GetNetatmoOAuth()

	if error != nil {
		logger.Error("Failed to get Netatmo historic data from Netatmo API, access token not working")
		result <- HistoricResult{History: models.NetatmoHistory{}, Error: error}
	}

	deviceId, deviceIdExists := os.LookupEnv("NETATMO_DEVICEID")

	if !deviceIdExists {
		logger.Error("Netatmo deviceid not set up in .env file! Can't fetch data from Netatmo API!")
		result <- HistoricResult{History: models.NetatmoHistory{}, Error: error}
	}

	apiUrl := utilities.BuildNetatmoMeasureUrl(token.AccessToken, deviceId, "", start.Unix(), end.Unix()).String()

	response, error := http.Get(apiUrl)

	if error != nil {
		logger.WithError(error).Error("Failed to get _history_ data from Netatmo API")
		result <- HistoricResult{History: models.NetatmoHistory{}, Error: error}
	}

	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		logger.Error("Failed to parse body from Netatmo _history_ API")
		result <- HistoricResult{History: models.NetatmoHistory{}, Error: error}
	}

	var historyData models.NetatmoHistory

	json.Unmarshal(body, &historyData)

	// TODO: Use mapper here to return our data model instead..

	result <- HistoricResult{History: historyData, Error: nil}
}

func (rest *rest) GetCurrent(result chan CurrentResult) {
	token, error := rest.credentials.GetNetatmoOAuth()

	if error != nil {
		logger.Error("Failed to get Netatmo current data from Netatmo API, access token not working")
		result <- CurrentResult{Current: models.NetatmoCurrent{}, Error: error}
	}

	deviceId, deviceIdExists := os.LookupEnv("NETATMO_DEVICEID")

	if !deviceIdExists {
		logger.Error("Netatmo deviceid not set up in .env file! Can't fetch data from Netatmo API!")
		result <- CurrentResult{Current: models.NetatmoCurrent{}, Error: error}
	}

	apiUrl := utilities.BuildStationUrl(token.AccessToken, deviceId).String()

	response, error := http.Get(apiUrl)

	if error != nil {
		logger.WithError(error).Error("Failed to get _current_ data from Netatmo API")
		result <- CurrentResult{Current: models.NetatmoCurrent{}, Error: error}
	}

	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		logger.Error("Failed to parse body from Netatmo _current_ API")
		result <- CurrentResult{Current: models.NetatmoCurrent{}, Error: error}
	}

	var currentData models.NetatmoCurrent

	json.Unmarshal(body, &currentData)

	// TODO: Use mapper here to return our data model instead..

	result <- CurrentResult{Current: currentData, Error: nil}
}
