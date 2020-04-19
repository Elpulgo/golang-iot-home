package wunderlist

import (
	"iot-home/models"
)

type Service interface {
	GetData(chan WunderlistResult) ([]models.WunderlistDto, error)
}

type service struct {
	repository RestService
}

type WunderlistResult struct {
	Lists []models.WunderlistDto
	Error error
}

func New(repository RestService) Service {
	return &service{repository: repository}
}

func (service *service) GetData() ([]models.WunderlistDto, error) {
	channel := make(chan WunderlistResult)
	go service.GetData(channel)
	response := <-channel

	return response.Lists, response.Error
}
