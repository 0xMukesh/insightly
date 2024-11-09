package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

type Llm string

var (
	Gemini  Llm = "gemini"
	Mistral     = "mistral"
	Llama       = "llama"
	Claude      = "claude"
)

type LlmConfig struct {
	Name   Llm    `json:"name"`
	ApiKey string `json:"api_key"`
}

type ConfigFile struct {
	Default Llm         `json:"default"`
	Llms    []LlmConfig `json:"llms"`
}

func GetConfigFilePath() string {
	homedir, _ := os.UserHomeDir()
	configFilePath := fmt.Sprintf("%s/.something.config.json", homedir)
	return configFilePath
}

func WriteToConfigFile(config ConfigFile) error {
	bytes, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(GetConfigFilePath(), bytes, 0644); err != nil {
		return err
	}

	return nil
}

func ReadConfigFile() (ConfigFile, error) {
	bytes, err := os.ReadFile(GetConfigFilePath())
	if err != nil {
		return ConfigFile{}, err
	}

	var configFile ConfigFile

	if err := json.Unmarshal(bytes, &configFile); err != nil {
		return ConfigFile{}, err
	}

	return configFile, nil
}
