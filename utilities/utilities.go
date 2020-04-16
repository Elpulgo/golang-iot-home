package utilities

import (
	"net/url"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

const NetatmoBaseUrl string = "https://api.netatmo.com"
const WunderlistBaseUrl string = "https://a.wunderlist.com/api/v1"

func BuildOauthTokenUrl() *url.URL {
	oauthUrl, error := url.Parse(NetatmoBaseUrl + "/oauth2/token")

	if error != nil {
		logger.WithError(error).Error("Failed to build Netatmo oauth2 token url")
	}

	return oauthUrl
}

func BuildNetatmoMeasureUrl(
	accessToken string,
	deviceId string,
	moduleId string,
	start int64,
	end int64,
) *url.URL {
	measureUrl, error := url.Parse(
		NetatmoBaseUrl +
			"/api/getmeasure" +
			"?access_token=" + accessToken +
			"&device_id=" + deviceId +
			"&module_id=" + moduleId +
			"&date_begin=" + strconv.FormatInt(start, 10) +
			"&date_end=" + strconv.FormatInt(end, 10) +
			"&scale=max" +
			"&type=temperature,humidity")

	if error != nil {
		logger.WithError(error).Error("Failed to build Netatmo measure url")
	}

	return measureUrl
}
func BuildStationUrl(accessToken string, deviceId string) *url.URL {
	stationUrl, error := url.Parse(
		NetatmoBaseUrl +
			"/api/getstationsdata" +
			"?access_token=" + accessToken +
			"&device_id=" + deviceId)

	if error != nil {
		logger.WithError(error).Error("Failed to build Netatmo station url")
	}
	return stationUrl
}

func BuildListsUrl(accessToken string, clientId string) *url.URL {
	listUrl, error := url.Parse(
		WunderlistBaseUrl +
			"/lists" +
			"?access_token=" + accessToken +
			"&client_id=" + clientId)

	if error != nil {
		logger.WithError(error).Error("Failed to build Wunderlist list url")
	}
	return listUrl
}

func BuildTasksUrl(accessToken string, clientId string, listId int64) *url.URL {
	tasksUrl, error := url.Parse(
		WunderlistBaseUrl +
			"/tasks" +
			"?access_token=" + accessToken +
			"&client_id=" + clientId +
			"&list_id=" + strconv.FormatInt(listId, 10))

	if error != nil {
		logger.WithError(error).Error("Failed to build Wunderlist tasks url")
	}
	return tasksUrl
}
