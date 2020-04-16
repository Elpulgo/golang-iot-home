package models

import "time"

type NetatmoHistory struct {
	Steps []Step `json:"body"`
	Name  string
}

type Step struct {
	Start    int64       `json:"beg_time"`
	Duration int         `json:"step_time"`
	Values   [][]float32 `json:"value"`
}

// Dto

type NetatmoSerieDto struct {
	Name   string
	Type   seriesType
	Values []netatmoValueDto
}

type netatmoValueDto struct {
	Value     float32
	Timestamp time.Time
}

type seriesType string

const (
	Temperature seriesType = "Temperature"
	Humidity    seriesType = "Humidity"
)

func (history NetatmoHistory) MapToDto(name string) []NetatmoSerieDto {
	var dtos []NetatmoSerieDto

	var valuesTemp []netatmoValueDto

	for _, step := range history.Steps {
		for _, value := range step.Values {
			valuesTemp = append(valuesTemp, netatmoValueDto{
				Timestamp: time.Unix(step.Start, 10),
				Value:     value[0],
			})
		}
	}

	dtos = append(dtos, NetatmoSerieDto{
		Name:   name + " Temperature",
		Type:   Temperature,
		Values: valuesTemp,
	})

	var valuesHumidity []netatmoValueDto

	for _, step := range history.Steps {
		for _, value := range step.Values {
			valuesHumidity = append(valuesHumidity, netatmoValueDto{
				Timestamp: time.Unix(step.Start, 10),
				Value:     value[len(value)-1],
			})
		}
	}

	dtos = append(dtos, NetatmoSerieDto{
		Name:   name + " Humidity",
		Type:   Humidity,
		Values: valuesHumidity,
	})

	return dtos
}
