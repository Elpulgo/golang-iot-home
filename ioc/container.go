package ioc

import (
	"iot-home/credentials"
	"iot-home/netatmo"

	"github.com/golobby/container"
)

func Setup() {
	container.Singleton(func() credentials.CredentialsService {
		return credentials.New()
	})

	container.Transient(func() netatmo.RestService {
		var credentials credentials.CredentialsService
		container.Make(&credentials)
		return netatmo.NewRestService(credentials)
	})

	container.Transient(func() netatmo.Service {
		var rest netatmo.RestService
		container.Make(&rest)
		return netatmo.NewService(rest)
	})
}
