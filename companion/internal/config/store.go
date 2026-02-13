package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStore struct {
	configPath string `json:"-"` // omit from json
	Port       int    `json:"port"`
}

func NewConfigStore(configPath string, port int) (*ConfigStore, error) {
	return &ConfigStore{
		configPath: configPath,
		Port:       port,
	}, nil
}

func CreateConfig(configPath string, port int) (*ConfigStore, error) {
	cs := &ConfigStore{
		configPath: configPath,
		Port:       port,
	}

	if err := cs.Save(); err != nil {
		return nil, err
	}

	return cs, nil
}

func LoadConfig(configPath string) (*ConfigStore, error) {
	cs := &ConfigStore{
		configPath: configPath,
	}

	if err := cs.Load(); err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *ConfigStore) Save() error {
	data, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(cs.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}

func (cs *ConfigStore) Load() error {
	data, err := os.ReadFile(cs.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	if err := json.Unmarshal(data, cs); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}

func (cs *ConfigStore) GetConfigPath() string {
	return cs.configPath
}

func (cs *ConfigStore) GetPort() int {
	return cs.Port
}

func ConfigExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return err == nil
}
