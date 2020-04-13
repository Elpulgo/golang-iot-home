package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"iot-home/netatmo"
	"iot-home/utilities"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	logger "github.com/sirupsen/logrus"
)

var lock = &sync.Mutex{}

var netatmoOAuth *netatmo.NetatmoOAuth
var hueAppKey *string

var DeviceName string = "b5c92462-aede-47de"
var AppName string = "IoTHomeDashboard"

type CredentialsService interface {
	GetHueCredentials() (appKey string, appName string, deviceName string)
	TryPersistHueAppKey(appKey string) bool
	GetWunderlistCredentials() (accessToken string, clientId string)
	GetNetatmoOAuth() (netatmo.NetatmoOAuth, error)
}

type credentialsService struct {
	service CredentialsService
}

func New() CredentialsService {
	return new(credentialsService)
}

func (credentialsService *credentialsService) GetHueCredentials() (appKey string, appName string, deviceName string) {
	appKey, error := loadHueAppKey()

	if error != nil {
		logger.WithError(error).Error("Failed to load app key")
		return "", "", ""
	}

	return appKey, AppName, DeviceName
}

func (credentialsService *credentialsService) TryPersistHueAppKey(appKey string) bool {
	if appKey == "" {
		return false
	}

	error := ioutil.WriteFile(hueAppKeyPath(), []byte(appKey), 0644)

	if error != nil {
		logger.WithError(error).Error("Failed to persist hueapp key")
		return false
	}

	lock.Lock()
	defer lock.Unlock()

	hueAppKey = &appKey

	return true
}

func (credentialsService *credentialsService) GetWunderlistCredentials() (accessToken string, clientId string) {
	clientId, clientIdExists := os.LookupEnv("WUNDERLIST_CLIENTID")
	accessToken, accessTokenExists := os.LookupEnv("WUNDERLIST_CLIENTID")

	if !clientIdExists {
		logger.Error("Wunderlist credentials is missing, ensure clientId exists in .env file!")
		return "", ""
	}

	if !accessTokenExists {
		logger.Error("Wunderlist credentials is missing, ensure accesstoken exists in .env file!")
		return "", ""
	}

	return accessToken, clientId
}

func (credentialsService *credentialsService) GetNetatmoOAuth() (netatmo.NetatmoOAuth, error) {

	if netatmoTokenAlreadyValid() {
		return *netatmoOAuth, nil
	}

	clientId, clientIdExists := os.LookupEnv("NETATMO_CLIENTID")
	clientSecret, clientSecretExists := os.LookupEnv("NETATMO_CLIENTSECRET")
	userName, userNameExists := os.LookupEnv("NETATMO_USERNAME")
	password, passwordExists := os.LookupEnv("NETATMO_PASSWORD")

	if !clientIdExists || !clientSecretExists || !userNameExists || !passwordExists {
		logger.Error("Netatmo credentials is missing, ensure clientid, clientsecret, username & password exsists in .env file!")
		return *netatmoOAuth, errors.New("Credentials for Netatmo does not exists in .env file!")
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
		logger.WithError(error).Error("Failed to get Netatmo OAuthToken")
		return *netatmoOAuth, errors.New("Failed to get Netatmo OAuth token from response!")
	}

	defer response.Body.Close()

	token, error := getNetatmoToken(response.Body)

	if error != nil {
		return *netatmoOAuth, errors.New("Failed to parse content from Netatmo OAuth request!")
	}

	lock.Lock()
	defer lock.Unlock()

	if netatmoOAuth == nil {
		netatmoOAuth = &token
	}

	return token, nil
}

func netatmoTokenAlreadyValid() bool {
	if netatmoOAuth != nil && !netatmoOAuth.HasExpired() {
		return true
	}

	return false
}

func getNetatmoToken(reader io.ReadCloser) (netatmo.NetatmoOAuth, error) {
	token := new(netatmo.NetatmoOAuth)
	err := json.NewDecoder(reader).Decode(token)

	if err != nil {
		logger.WithError(err).Error("Failed to read OAuth token from Netatmo request")
		return *token, err
	}

	return *token, nil
}

func loadHueAppKey() (string, error) {
	if hueAppKey != nil {
		return *hueAppKey, nil
	}

	keyPath := hueAppKeyPath()

	fileInfo, error := os.Stat(keyPath)
	_ = fileInfo

	if error != nil {
		return "", error
	}

	fileContent, error := ioutil.ReadFile(keyPath)

	if error != nil {
		return "", error
	}

	lock.Lock()
	defer lock.Unlock()

	if hueAppKey == nil {
		key := string(fileContent)
		hueAppKey = &key
	}

	return string(fileContent), nil
}

func hueAppKeyPath() string {
	currentDir, error := os.Getwd()

	if error != nil {
		panic(fmt.Sprintf("Failed to find current directory .."))
	}

	fileInfo, error := os.Stat(path.Join(currentDir, "settings", "hueappkey.dat"))
	_ = fileInfo

	if os.IsNotExist(error) {
		if errCreateDir := os.MkdirAll(path.Join(currentDir, "settings"), 0755); errCreateDir != nil {
			logger.WithError(errCreateDir).Error("Failed to create directory for hueappkey to be stored")
			return ""
		}
	}

	return path.Join(currentDir, "settings", "hueappkey.dat")
}
