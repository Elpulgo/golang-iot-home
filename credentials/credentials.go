package credentials

import (
	"fmt"
	"iot-home/logger"
	"iot-home/utilities"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

	u, _ := url.ParseRequestURI(apiUrl.EscapedPath())
	// u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	fmt.Println(urlStr)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl.EscapedPath(), strings.NewReader(data.Encode())) // URL-encoded payload
	// r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, error := client.Do(r)
	if error != nil {
		logger.Error(fmt.Sprintf("Failed to get Netatmo OAuthToken %s", error.Error()))
		return
	}
	fmt.Println(resp.Body)
}
