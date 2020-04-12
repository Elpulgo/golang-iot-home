package ioc

import (
	"iot-home/credentials"
	"iot-home/netatmo"
	"iot-home/netatmoRest"

	"github.com/golobby/container"
)

func Setup() {
	container.Singleton(func() credentials.CredentialsService {
		return credentials.New()
	})

	container.Transient(func() netatmoRest.RestService {
		return netatmoRest.New()
	})

	container.Transient(func() netatmo.Service {
		var rest netatmoRest.RestService
		container.Make(&rest)
		return netatmo.New(rest)
	})
}
