package main

import (
	"log"

	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"
)

type SettingsInfo struct {
	Active          bool
	IntervalMinutes int
}

const settingsFilename = "settings.yaml"

var defaultSettings = SettingsInfo{true, 60}
var configFolder = configdir.New("ekevoo", "cuckoo").QueryFolders(configdir.Global)[0]
var settings = NewSettings()

func NewSettings() SettingsInfo {
	payload, err := configFolder.ReadFile(settingsFilename)
	if err != nil {
		log.Print(err)
		return defaultSettings
	}

	var result SettingsInfo
	if err := yaml.Unmarshal(payload, &result); err != nil {
		log.Print(err)
		return defaultSettings
	}

	return result
}

func (settings SettingsInfo) Save() error {
	payload, err := yaml.Marshal(&settings)
	if err != nil {
		return err
	}
	return configFolder.WriteFile(settingsFilename, payload)
}
