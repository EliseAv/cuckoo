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

var defaults = SettingsInfo{true, 60}
var settings = NewSettings()

var configFolder = configdir.New("ekevoo", "cuckoo").QueryFolders(configdir.Local)[0]

func NewSettings() SettingsInfo {
	configFolder.MkdirAll()
	payload, err := configFolder.ReadFile("settings.json")
	if err != nil {
		log.Print(err)
		return defaults
	}

	var result SettingsInfo
	if err := yaml.Unmarshal(payload, &result); err != nil {
		log.Print(err)
		return defaults
	}

	return result
}

func (settings SettingsInfo) Save() error {
	payload, err := yaml.Marshal(&settings)
	if err != nil {
		return err
	}
	return configFolder.WriteFile("settings.json", payload)
}
