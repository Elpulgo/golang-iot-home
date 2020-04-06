package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"iot-home/logger"
	"iot-home/netatmo"
	"iot-home/utilities"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

var lock = &sync.Mutex{}

var netatmoOAuth *netatmo.NetatmoOAuth
var hueAppKey *string

const deviceName string = "b5c92462-aede-47de"
const appName string = "IoTHomeDashboard"

func GetHueCredentials() (appKey string, appName string, deviceName string) {
	appKey, error := loadAppKey()

	if error != nil {
		logger.Error(fmt.Sprintf("Failed to load app key %s", error.Error()))
		return "", "", ""
	}

	return appKey, appName, deviceName
}

func TryPersistHueAppKey(appKey string) bool {
	if appKey == "" {
		return false
	}

	error := ioutil.WriteFile(hueAppKeyPath(), []byte(appKey), 0644)

	if error != nil {
		logger.Error(fmt.Sprintf("Failed to persist hue app key %s", error.Error()))
		return false
	}

	lock.Lock()
	defer lock.Unlock()

	hueAppKey = &appKey

	return true
}

func GetWunderlistCredentials() {

}

func GetNetatmoOAuth() (*netatmo.NetatmoOAuth, error) {

	if tokenAlreadyValid() {
		return netatmoOAuth, nil
	}

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

	lock.Lock()
	defer lock.Unlock()

	if netatmoOAuth == nil {
		netatmoOAuth = &token
	}

	return &token, nil
}

func tokenAlreadyValid() bool {
	if netatmoOAuth != nil && !netatmoOAuth.HasExpired() {
		return true
	}

	return false
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

func loadAppKey() (string, error) {
	if hueAppKey != nil {
		return *hueAppKey, nil
	}

	fileInfo, error := os.Stat(hueAppKeyPath())

	if os.IsNotExist(error) {
		return "", error
	}

	fileContent, error := ioutil.ReadFile(fileInfo.Name())

	if error != nil {
		return "", error
	}

	return string(fileContent), nil

}

func hueAppKeyPath() string {
	currentDir, error := os.Getwd()

	if error != nil {
		panic(fmt.Sprintf("Failed to find current directory .."))
	}

	return path.Join(currentDir, "settings", "hueappkey.dat")
}
