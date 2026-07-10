package config

import (
	"encoding/json"
	"log"
	"os"
)

func CreateDefaultConfigFile() error {
	defaultConfig := Config{
		LibraryPath: "library",
		Targets: []Target{
			{
				Name:   "Original Size",
				Width:  0,
				Height: 0,
				Color:  true,
				AutoRotate: false,
			},
		},
	}
	return SaveConfigFile(&defaultConfig)
}

func SaveConfigFile(config *Config) error {
	// marshal the config struct to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// write the JSON data to the config file
	err = os.WriteFile(APPCONFIG_FILE, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetConfigFile() (*Config, error) {
	// check if the config file exists if not create it
	if _, err := os.Stat(APPCONFIG_FILE); os.IsNotExist(err) {
		err := CreateDefaultConfigFile()
		if err != nil {
			return nil, err
		}
	}

	// read the config file
	data, err := os.ReadFile(APPCONFIG_FILE)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfigFile() error {
	cfg, err := GetConfigFile()
	if err != nil {
		return err
	}

	Configuration = *cfg
	log.Printf("Current configuration: %+v", Configuration)
	return nil
}
