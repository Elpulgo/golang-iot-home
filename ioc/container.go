package ioc

import (
	"iot-home/credentials"
	"iot-home/hue"
	"iot-home/netatmo"
	"iot-home/wunderlist"

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

	container.Transient(func() wunderlist.RestService {
		var credentials credentials.CredentialsService
		container.Make(&credentials)
		return wunderlist.NewRestService(credentials)
	})

	container.Transient(func() wunderlist.Service {
		var rest wunderlist.RestService
		container.Make(&rest)
		return wunderlist.NewWunderlistService(rest)
	})

	container.Transient(func() hue.Registry {
		var credentials credentials.CredentialsService
		container.Make(&credentials)
		return hue.NewRegistry(credentials)
	})

	container.Transient(func() hue.Service {
		var registry hue.Registry
		container.Make(&registry)
		return hue.NewHueService(registry)
	})
}
