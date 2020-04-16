package models

type NetatmoCurrent struct {
	CurrentData CurrentData `json:"body"`
}

type CurrentData struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	DashboardData DashboardData `json:"dashboard_data"`
	Modules       []Module      `json:"modules"`
}

type Module struct {
	DashboardData DashboardData `json:"dashboard_data"`
}

type DashboardData struct {
	Temperature      float32 `json:"temperature"`
	CO2              float32 `json:"co2"`
	Pressure         float32 `json:"pressure"`
	AbsolutePressure float32 `json:"absolutePressure"`
	MinTemp          float32 `json:"min_temp"`
	MaxTemp          float32 `json:"max_temp"`
	Humidity         int64   `json:"humidity"`
	TempTrend        string  `json:"temp_trend"`
	PressureTrend    string  `json:"pressure_trend"`
}

// Dto

type NetatmoCurrentDto struct {
	Name             string
	Temperature      float32
	CO2              float32
	Humidity         int64
	Pressure         float32
	AbsolutePressure float32
	MinTemp          float32
	MaxTemp          float32
	TempTrend        string
	PressureTrend    string
}

func (data CurrentData) MapToDto() []NetatmoCurrentDto {
	var dtos []NetatmoCurrentDto

	if len(data.Devices) < 1 {
		return dtos
	}

	indoorModule := data.Devices[0]

	dtos = append(dtos, NetatmoCurrentDto{
		Name:             "Indoor",
		Temperature:      indoorModule.DashboardData.Temperature,
		CO2:              indoorModule.DashboardData.CO2,
		Humidity:         indoorModule.DashboardData.Humidity,
		Pressure:         indoorModule.DashboardData.Pressure,
		AbsolutePressure: indoorModule.DashboardData.AbsolutePressure,
		MinTemp:          indoorModule.DashboardData.MinTemp,
		MaxTemp:          indoorModule.DashboardData.MaxTemp,
		TempTrend:        indoorModule.DashboardData.TempTrend,
		PressureTrend:    indoorModule.DashboardData.PressureTrend,
	})

	if len(data.Devices[0].Modules) < 1 {
		return dtos
	}

	outdoorModule := data.Devices[0].Modules[0]

	dtos = append(dtos, NetatmoCurrentDto{
		Name:             "Outdoor",
		Temperature:      outdoorModule.DashboardData.Temperature,
		CO2:              outdoorModule.DashboardData.CO2,
		Humidity:         outdoorModule.DashboardData.Humidity,
		Pressure:         outdoorModule.DashboardData.Pressure,
		AbsolutePressure: outdoorModule.DashboardData.AbsolutePressure,
		MinTemp:          outdoorModule.DashboardData.MinTemp,
		MaxTemp:          outdoorModule.DashboardData.MaxTemp,
		TempTrend:        outdoorModule.DashboardData.TempTrend,
		PressureTrend:    outdoorModule.DashboardData.PressureTrend,
	})

	return dtos
}
