package netatmo

type Rest struct {
}

func (rest *Rest) GetCurrent() (Current, error) {
	current := &Current{
		CurrentData: CurrentData{
			Devices: []Device{
				Device{
					DashboardData: DashboardData{
						Temperature:      32,
						CO2:              1,
						Pressure:         3,
						AbsolutePressure: 5,
						MinTemp:          4,
						MaxTemp:          5,
						Humidity:         30,
						TempTrend:        "Up",
						PressureTrend:    "Down",
					},
				},
			},
		},
	}

	return *current, nil
}
