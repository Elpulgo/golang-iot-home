package hue

import (
	"errors"
	"iot-home/credentials"
	"os"
	"sort"

	"github.com/amimof/huego"
	logger "github.com/sirupsen/logrus"
)

type Registry interface {
	Connect() (huego.Bridge, error)
}

type registry struct {
	credentials credentials.CredentialsService
}

func NewRegistry(credentials credentials.CredentialsService) Registry {
	return &registry{credentials: credentials}
}

func (registry *registry) Connect() (huego.Bridge, error) {

	bridge := huego.Bridge{}
	bridgeIp, bridgeIpExists := os.LookupEnv("PHILIPS_HUE_BRIDGEIP")

	if !bridgeIpExists {
		logger.Error("No Philips Hue bridge Ip exists in .env, can't connect to Hue bridge!")
		return bridge, errors.New("Failed to connect to Hue Brdige! .env variable for IP is missing!")
	}

	appKey, appName, deviceName := registry.credentials.GetHueCredentials()

	if appKey == "" {
		logger.Info("Failed to load app key for Philips Hue bridge, will try and register ...")
		bridge, succeded := tryRegister(appName, deviceName, bridgeIp)
		if !succeded {
			logger.Error("Failed to locate bridge, won't connect")
			return huego.Bridge{}, errors.New("Failed to locate bridge, won't connect")
		}

		logger.Info(bridge)
	}

	return huego.Bridge{}, nil
	// bridge, succeded := huego.New("192.168.1.59", "username")
}

func tryRegister(appName string, deviceName string, bridgeIp string) (huego.Bridge, bool) {
	bridges, error := huego.DiscoverAll()
	if error != nil {
		logger.WithError(error).Error("Failed to locate Hue bridges on the network")
		return huego.Bridge{}, false
	}

	index := sort.Search(len(bridges), func(i int) bool {
		return string(bridges[i].Host) >= bridgeIp
	})

	if index > 0 {
		return huego.Bridge{}, false
	}

	bridge := bridges[index]
	user, error := bridge.CreateUser(appName)
	if error != nil {
		logger.WithError(error).Error("Failed to create user for Hue bridge")
		return huego.Bridge{}, false
	}

	connectedBridge := bridge.Login(user)

	logger.Info(connectedBridge)

	return huego.Bridge{}, false
}
