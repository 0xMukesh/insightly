package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Llm string

var (
	Gemini  Llm = "gemini"
	Mistral Llm = "mistral"
	Llama   Llm = "llama"
	Claude  Llm = "claude"
	Chatgpt Llm = "chatgpt"
	Qwen    Llm = "qwen"
)

var ValidLlms = []string{string(Gemini), string(Mistral), string(Llama), string(Claude), string(Chatgpt), string(Qwen)}

type LlmConfig struct {
	Name   Llm    `json:"name" mapstructure:"name"`
	ApiKey string `json:"api_key" mapstructure:"api_key"`
}

type ConfigFile struct {
	Default Llm         `json:"default" mapstructure:"default"`
	Llms    []LlmConfig `json:"llms" mapstructure:"llms"`
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

func DoesConfigFileExists() bool {
	if _, err := os.Stat(GetConfigFilePath()); err != nil {
		return false
	}

	return true
}

func GetLlmKey(llmName string) (string, error) {
	config, err := ReadConfigFile()
	if err != nil {
		return "", err
	}

	for i := range config.Llms {
		if config.Llms[i].Name == Llm(strings.ToLower(llmName)) {
			return config.Llms[i].ApiKey, nil
		}
	}

	return "", errors.New("invalid llm")
}
