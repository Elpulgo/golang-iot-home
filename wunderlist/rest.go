package wunderlist

import (
	"encoding/json"
	"io/ioutil"
	"iot-home/credentials"
	"iot-home/models"
	"iot-home/utilities"
	"net/http"
	"os"
	"strconv"
	"sync"

	logger "github.com/sirupsen/logrus"
)

type RestService interface {
	GetData(result chan WunderlistResult)
}

type rest struct {
	service     RestService
	credentials credentials.CredentialsService
}

func NewRestService(credentials credentials.CredentialsService) RestService {
	return &rest{credentials: credentials}
}

func (rest *rest) GetData(result chan WunderlistResult) {
	accessToken, clientId := rest.credentials.GetWunderlistCredentials()

	if accessToken == "" || clientId == "" {
		logger.Error("Failed to get Wunderlist credentials")
		result <- WunderlistResult{}
		return
	}

	_, firstListExists := os.LookupEnv("WUNDERLIST_LISTFIRST")

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

	if error != nil {
		logger.WithError(error).Error("Failed to read body from Wunderlist API list response")
		result <- WunderlistResult{}
		return
	}

	var listsData []models.WunderlistListData
	json.Unmarshal(body, &listsData)

	wunderlistDtos, error := getTasks(listsData, accessToken, clientId)

	result <- WunderlistResult{Lists: wunderlistDtos, Error: nil}
}

func getTasks(
	listsData []models.WunderlistListData,
	accessToken string,
	clientId string) ([]models.WunderlistDto, error) {

	lists := getLists()

	wunderlistData := filterLists(listsData, lists)

	var dtos []models.WunderlistDto
	var waitGroup sync.WaitGroup

	queue := make(chan models.WunderlistListData, 1)

	waitGroup.Add(len(wunderlistData))

	for _, data := range wunderlistData {
		go func(data models.WunderlistListData) {
			queue <- data
		}(data)
	}

	go func() {
		for data := range queue {
			apiUrl := utilities.BuildWunderlistTasksUrl(accessToken, clientId, data.Id).String()
			response, error := http.Get(apiUrl)

			if error != nil {
				logger.WithError(error).Error("Failed to get _tasks_ data from Wunderlist API")
				defer response.Body.Close()
				continue
			}

			defer response.Body.Close()
			body, error := ioutil.ReadAll(response.Body)

			if error != nil {
				logger.WithError(error).WithField("listId", strconv.FormatInt(data.Id, 10)).Error("Failed to read response from Wunderlist API for tasks")
				continue
			}

			var tasksData []models.WunderlistTaskData
			json.Unmarshal(body, &tasksData)
			tasksDto := models.MapToDto(tasksData, data.Name)

			dtos = append(dtos, tasksDto)
			waitGroup.Done()
		}
	}()

	waitGroup.Wait()

	return dtos, nil
}

func getLists() []string {
	var lists []string

	firstList, firstListExists := os.LookupEnv("WUNDERLIST_LISTFIRST")

	if !firstListExists {
		logger.Error("No list exists in .env, can't fetch Wunderlist data!")
		return lists
	}

	lists = append(lists, firstList)

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
	filtered := make(map[string]struct{}, len(selectedLists))
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
