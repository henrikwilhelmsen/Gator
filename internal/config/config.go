package config

import (
	"encoding/json"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
	// Environment variable for overriding the config filepath
	ConfigFilePathEnvVar = "GATOR_CONFIG_FILEPATH"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Get the config filepath. Uses ConfigFilePathEnvVar if set, otherwise the default path
func getConfigFilePath() (string, error) {
	if envPath := os.Getenv(ConfigFilePathEnvVar); envPath != "" {
		return envPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

// Read the given filepath to a Config struct
func readFromFile(filepath string) (Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// Write the given config to the given filepath
func writeToFile(filepath string, cfg *Config) error {
	data, err := json.Marshal(*cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0o666)
	if err != nil {
		return err
	}

	return nil
}

// Reads the config file and returns the data as a Config struct
func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	return readFromFile(configFilePath)
}

// Sets the CurrentUserName to the given string and writes the config to file
func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return writeToFile(configFilePath, cfg)
}
