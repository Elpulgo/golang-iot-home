package netatmo

import (
	"iot-home/models"
	"time"
)

type Service interface {
	GetCurrent() (models.NetatmoCurrent, error)
	GetHistory(start time.Time, end time.Time) (models.NetatmoHistory, error)
}

type service struct {
	repository RestService
}

type CurrentResult struct {
	Current models.NetatmoCurrent
	Error   error
}

type HistoricResult struct {
	History models.NetatmoHistory
	Error   error
}

func NewService(repository RestService) Service {
	return &service{repository: repository}
}

func (service *service) GetCurrent() (models.NetatmoCurrent, error) {
	channel := make(chan CurrentResult)

	go service.repository.GetCurrent(channel)

	response := <-channel

	return response.Current, response.Error
}

func (service *service) GetHistory(start time.Time, end time.Time) (models.NetatmoHistory, error) {
	channel := make(chan HistoricResult)

	go service.repository.GetHistory(start, end, channel)

	response := <-channel
	return response.History, response.Error
}
