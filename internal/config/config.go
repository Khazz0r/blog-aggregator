package config

import (
	"os"
	"encoding/json"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read and unmarshal the .gatorconfig.json data into a Config struct
func Read() (Config, error) {
	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Write .gatorconfig.json file with the Config struct gotten from SetUser()
func write(cfg Config) error {
	configFile, _ := getConfigFilePath()
	
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + configFileName, nil
}
