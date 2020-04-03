package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iot-home/logger"
	"iot-home/netatmo"
	"iot-home/utilities"
	"net/http"
	"net/url"
	"os"
)

func GetHueCredentials() {

}

func GetWunderlistCredentials() {

}

func GetNetatmoOAuth() (*netatmo.NetatmoOAuth, error) {
	clientId, clientIdExists := os.LookupEnv("NETATMO_CLIENTID")
	clientSecret, clientSecretExists := os.LookupEnv("NETATMO_CLIENTSECRET")
	userName, userNameExists := os.LookupEnv("NETATMO_USERNAME")
	password, passwordExists := os.LookupEnv("NETATMO_PASSWORD")

	if !clientIdExists || !clientSecretExists || !userNameExists || !passwordExists {
		logger.Error("Netatmo credentials is missing, ensure clientid, clientsecret, username & password exsists in .env file!")
		return nil, errors.New("Credentials for Netatmo does not exists in .env file!")
	}

	apiUrl := utilities.BuildOauthTokenUrl().String()

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("username", userName)
	data.Set("password", password)

	response, error := http.PostForm(apiUrl, data)

	if error != nil {
		logger.Error(fmt.Sprintf("Failed to get Netatmo OAuthToken %s", error.Error()))
		return nil, errors.New("Failed to get Netatmo OAuth token from response!")
	}

	defer response.Body.Close()

	token, error := getToken(response.Body)

	if error != nil {
		return nil, errors.New("Failed to parse content from Netatmo OAuth request!")
	}

	return &token, nil
}

func getToken(reader io.ReadCloser) (netatmo.NetatmoOAuth, error) {
	token := new(netatmo.NetatmoOAuth)
	err := json.NewDecoder(reader).Decode(token)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read OAuth token from Netatmo request %s", err.Error()))
		return *token, err
	}

	return *token, nil
}
