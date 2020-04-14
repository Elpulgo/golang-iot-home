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
