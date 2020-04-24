package models

import (
	"github.com/amimof/huego"
)

type HueDto struct {
	Lights []HueLight
}

type HueLight struct {
	Name       string
	HexColor   string
	Saturation int64
	On         bool
}

func MapLightToDto(huegoLight huego.Light) HueLight {
	return HueLight{
		Name:       huegoLight.Name,
		On:         huegoLight.IsOn(),
		HexColor:   "FFF",
		Saturation: 100,
	}
}
