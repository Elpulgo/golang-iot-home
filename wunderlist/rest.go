package wunderlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"iot-home/credentials"
	"iot-home/models"
	"iot-home/utilities"
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"
)

type RestService interface {
	GetData(result chan WunderlistResult)
}

type rest struct {
	service     RestService
	credentials credentials.CredentialsService
}

func New(credentials credentials.CredentialsService) RestService {
	return &rest{credentials: credentials}
}

func (rest *rest) GetData(result chan WunderlistResult) {
	accessToken, clientId := rest.credentials.GetWunderlistCredentials()

	if accessToken == "" || clientId == "" {
		logger.Error("Failed to get Wunderlist credentials")
		result <- WunderlistResult{}
		return
	}

	firstList, firstListExists := os.LookupEnv("WUNDERLIST_LISTFIRST")
	secondList, secondListExists := os.LookupEnv("WUNDERLIST_LISTSECOND")

	if !firstListExists {
		logger.Error("No list exists in .env, can't fetch Wunderlist data!")
		result <- WunderlistResult{}
		return
	}

	apiUrl := utilities.BuildWunderlistListsUrl(accessToken, clientId).String()
	response, error := http.Get(apiUrl)

	if error != nil {
		logger.WithError(error).Error("Failed to get _lists_ data from Wunderlist API")
		result <- WunderlistResult{}
		return
	}

	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)

	var listsData []models.WunderlistListData

	json.Unmarshal(body, &listsData)
}

func getTasks(listsData []models.WunderlistListData) (models.WunderlistDto, error) {
	lists := getLists()

	wunderlistData := filterLists(listsData, lists)

	fmt.Println(wunderlistData)
}

func getLists() []string {
	var lists []string

	firstList, firstListExists := os.LookupEnv("WUNDERLIST_LISTFIRST")

	if !firstListExists {
		logger.Error("No list exists in .env, can't fetch Wunderlist data!")
		return lists
	}

	secondList, secondListExists := os.LookupEnv("WUNDERLIST_LISTSECOND")
	thirdList, thirdListExists := os.LookupEnv("WUNDERLIST_LISTTHIRD")
	fourthList, fourthListExists := os.LookupEnv("WUNDERLIST_LISTFOURTH")
	fifthList, fifthListExists := os.LookupEnv("WUNDERLIST_LISTFIFTH")

	if secondListExists {
		lists = append(lists, secondList)
	}
	if thirdListExists {
		lists = append(lists, thirdList)
	}
	if fourthListExists {
		lists = append(lists, fourthList)
	}
	if fifthListExists {
		lists = append(lists, fifthList)
	}

	return lists
}

func filterLists(listsFromResponse []models.WunderlistListData, selectedLists []string) (out []models.WunderlistListData) {
	filtered := make(map[models.WunderlistListData]struct{}, len(selectedLists))
	for _, data := range selectedLists {
		filtered[data] = struct{}{}
	}
	for _, data := range listsFromResponse {
		if _, ok := filtered[data.Name]; ok {
			out = append(out, data)
		}
	}
	return
}
