package hue

import (
	"iot-home/models"

	logger "github.com/sirupsen/logrus"
)

type Service interface {
	GetLights() (models.HueDto, error)
}

type service struct {
	registry Registry
}

func NewHueService(registry Registry) Service {
	return &service{registry: registry}
}

func (service *service) GetLights() (models.HueDto, error) {
	bridge, error := service.registry.Connect()
	if error != nil {
		logger.WithError(error).Error("Failed to connect to Hue Bridge!")
		return models.HueDto{}, error
	}

	lights, error := bridge.GetLights()
	if error != nil {
		logger.WithError(error).Error("Failed to get lights from Hue API!")
		return models.HueDto{}, error
	}

	var lightsDto []models.HueLight
	for _, light := range lights {
		lightsDto = append(lightsDto, models.MapLightToDto(light))
	}

	return models.HueDto{
		Lights: lightsDto,
	}, nil

}
