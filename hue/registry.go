package hue

import (
	"errors"
	"iot-home/credentials"
	"os"
	"sort"
	"strings"

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

	bridgeIp, bridgeIpExists := os.LookupEnv("PHILIPS_HUE_BRIDGEIP")

	if !bridgeIpExists {
		logger.Error("No Philips Hue bridge Ip exists in .env, can't connect to Hue bridge!")
		return huego.Bridge{}, errors.New("Failed to connect to Hue Brdige! .env variable for IP is missing!")
	}

	appKey, appName, deviceName := registry.credentials.GetHueCredentials()

	if appKey == "" {
		logger.Info("Failed to load app key for Philips Hue bridge, will try and register ...")
		bridge, succeded, error := registry.tryRegister(appName, deviceName, bridgeIp)
		if !succeded {
			return huego.Bridge{}, error
		}

		return bridge, nil
	}

	bridgeNew := huego.New(bridgeIp, appName)
	var connectedBridge huego.Bridge = *bridgeNew

	return connectedBridge, nil
}

func (registry *registry) tryRegister(appName string, deviceName string, bridgeIp string) (huego.Bridge, bool, error) {
	bridges, error := huego.DiscoverAll()
	if error != nil {
		logger.WithError(error).Error("Failed to locate Hue bridges on the network")
		return huego.Bridge{}, false, errors.New("Failed to locate Hue bridges on the network")
	}

	index := sort.Search(len(bridges), func(i int) bool {
		return string(bridges[i].Host) >= bridgeIp
	})

	if index > 0 {
		return huego.Bridge{}, false, errors.New("Failed to find bridge on network")
	}

	foundBridge := bridges[index]
	user, error := foundBridge.CreateUser(appName)
	if error != nil {
		if strings.Contains("ERROR 101", error.Error()) {
			logger.WithError(error).Error("Link button not pressed!")
			return huego.Bridge{}, false, errors.New("Link button not pressed!")
		}

		logger.WithError(error).Error("Failed to create user for Hue bridge")
		return huego.Bridge{}, false, errors.New("Failed to create user for Hue bridge")
	}

	bridge := foundBridge.Login(user)
	var connectedBridge huego.Bridge = *bridge
	logger.Info(connectedBridge)

	if !registry.credentials.TryPersistHueAppKey(connectedBridge.ID) {
		logger.Error("Failed to persist app key for Hue bridge ...")
		return huego.Bridge{}, false, errors.New("Failed to persist app key for Hue bridge!")
	}

	return connectedBridge, false, nil
}
