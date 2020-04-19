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
		result <- HistoricResult{History: []models.NetatmoSerieDto{}, Error: error}
		return
	}

	deviceId, deviceIdExists := os.LookupEnv("NETATMO_DEVICEID")

	if !deviceIdExists {
		logger.Error("Netatmo deviceid not set up in .env file! Can't fetch data from Netatmo API!")
		result <- HistoricResult{History: []models.NetatmoSerieDto{}, Error: error}
		return
	}

	channelIndoor := make(chan HistoricResult)
	go getHistory(deviceId, token.AccessToken, start, end, "", "Indoor", channelIndoor)
	responseIndoor := <-channelIndoor

	if responseIndoor.Error != nil {
		logger.Error(responseIndoor.Error.Error())
	}

	moduleId, moduleIdExists := os.LookupEnv("NETATMO_OUTDOORMODULEID")

	if !moduleIdExists {
		logger.Error("Netatmo outdoor module not set up in .env file! Can't fetch data from Netatmo API!")
		result <- HistoricResult{History: responseIndoor.History, Error: nil}
		return
	}

	channelOutdoor := make(chan HistoricResult)
	go getHistory(deviceId, token.AccessToken, start, end, moduleId, "Outdoor", channelOutdoor)
	responseOutdoor := <-channelOutdoor

	if responseOutdoor.Error != nil {
		logger.Error(responseOutdoor.Error.Error())
	}

	series := []models.NetatmoSerieDto{}
	series = append(responseIndoor.History, responseOutdoor.History...)

	result <- HistoricResult{History: series, Error: nil}
}

func getHistory(
	deviceId string,
	token string,
	start time.Time,
	end time.Time,
	moduleId string,
	name string,
	result chan HistoricResult) {

	apiUrl := utilities.BuildNetatmoMeasureUrl(token, deviceId, moduleId, start.Unix(), end.Unix()).String()

	response, error := http.Get(apiUrl)

	if error != nil {
		logger.WithError(error).Error("Failed to get _history_ data from Netatmo API")
		result <- HistoricResult{History: []models.NetatmoSerieDto{}, Error: error}
		return
	}

	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		logger.Error("Failed to parse body from Netatmo _history_ API")
		result <- HistoricResult{History: []models.NetatmoSerieDto{}, Error: error}
		return
	}

	var historyData models.NetatmoHistory

	json.Unmarshal(body, &historyData)
	result <- HistoricResult{History: historyData.MapToDto(name), Error: nil}
}

func (rest *rest) GetCurrent(result chan CurrentResult) {
	token, error := rest.credentials.GetNetatmoOAuth()

	if error != nil {
		logger.Error("Failed to get Netatmo current data from Netatmo API, access token not working")
		result <- CurrentResult{Current: []models.NetatmoCurrentDto{}, Error: error}
		return
	}

	deviceId, deviceIdExists := os.LookupEnv("NETATMO_DEVICEID")

	if !deviceIdExists {
		logger.Error("Netatmo deviceid not set up in .env file! Can't fetch data from Netatmo API!")
		result <- CurrentResult{Current: []models.NetatmoCurrentDto{}, Error: error}
		return
	}

	apiUrl := utilities.BuildStationUrl(token.AccessToken, deviceId).String()

	response, error := http.Get(apiUrl)

	if error != nil {
		logger.WithError(error).Error("Failed to get _current_ data from Netatmo API")
		result <- CurrentResult{Current: []models.NetatmoCurrentDto{}, Error: error}
		return
	}

	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		logger.Error("Failed to parse body from Netatmo _current_ API")
		result <- CurrentResult{Current: []models.NetatmoCurrentDto{}, Error: error}
		return
	}

	var currentData models.NetatmoCurrent

	json.Unmarshal(body, &currentData)

	result <- CurrentResult{Current: currentData.CurrentData.MapToDto(), Error: nil}
}
