package credentials

import (
	"fmt"
	"iot-home/logger"
	"iot-home/utilities"
	"net/http"
	"net/url"
	"os"
)

func GetHueCredentials() {

}

func GetWunderlistCredentials() {

}

func GetNetatmoOAuth() {
	clientId, clientIdExists := os.LookupEnv("NETATMO_CLIENTID")
	clientSecret, clientSecretExists := os.LookupEnv("NETATMO_CLIENTSECRET")
	userName, userNameExists := os.LookupEnv("NETATMO_USERNAME")
	password, passwordExists := os.LookupEnv("NETATMO_PASSWORD")

	if !clientIdExists || !clientSecretExists || !userNameExists || !passwordExists {
		logger.Error("Netatmo credentials is missing, ensure clientid, clientsecret, username & password exsists in .env file!")
		return
	}

	apiUrl := utilities.BuildOauthTokenUrl()

	fmt.Println(apiUrl)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("username", userName)
	data.Set("password", password)

	fmt.Println(data)
	fmt.Println(apiUrl.String())

	resp, error := http.PostForm(apiUrl.String(), data)

	if error != nil {
		logger.Error(fmt.Sprintf("Failed to get Netatmo OAuthToken %s", error.Error()))
		return
	}
	fmt.Println(resp.Body)
}
