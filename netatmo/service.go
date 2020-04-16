package netatmo

import (
	"iot-home/models"
	"time"
)

type Service interface {
	GetCurrent() ([]models.NetatmoCurrentDto, error)
	GetHistory(start time.Time, end time.Time) ([]models.NetatmoSerieDto, error)
}

type service struct {
	repository RestService
}

type CurrentResult struct {
	Current []models.NetatmoCurrentDto
	Error   error
}

type HistoricResult struct {
	History []models.NetatmoSerieDto
	Error   error
}

func NewService(repository RestService) Service {
	return &service{repository: repository}
}

func (service *service) GetCurrent() ([]models.NetatmoCurrentDto, error) {
	channel := make(chan CurrentResult)

	go service.repository.GetCurrent(channel)

	response := <-channel

	return response.Current, response.Error
}

func (service *service) GetHistory(start time.Time, end time.Time) ([]models.NetatmoSerieDto, error) {
	channel := make(chan HistoricResult)

	go service.repository.GetHistory(start, end, channel)

	response := <-channel
	return response.History, response.Error
}
