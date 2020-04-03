package utilities

import (
	"fmt"
	"iot-home/logger"
	"net/url"
)

const NetatmoBaseUrl string = "https://api.netatmo.com"

func BuildOauthTokenUrl() *url.URL {
	oauthUrl, error := url.Parse(NetatmoBaseUrl + "/oauth2/token")

	if error != nil {
		logger.Error(fmt.Sprintf("Failed to build Netatmo oauth2 token url, %s", error.Error()))
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
			"&date_begin=" + fmt.Sprintf("%s", string(start)) +
			"&date_end=" + fmt.Sprintf("%s", string(end)) +
			"&scale=max" +
			"&type=temperature,humidity")

	if error != nil {
		logger.Error(fmt.Sprintf("Failed to build Netatmo measure url, %s" + error.Error()))
	}

	return measureUrl
}
