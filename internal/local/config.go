package local

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"cours-go/internal/domain"
)

const appDirectory = "reservation-salles"

func LoadConfig() (domain.ClientConfig, error) {
	path, err := configFilePath("config.json")
	if err != nil {
		return domain.ClientConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := DefaultConfig()
			return cfg, SaveConfig(cfg)
		}
		return domain.ClientConfig{}, err
	}

	var cfg domain.ClientConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return domain.ClientConfig{}, err
	}

	return NormalizeConfig(cfg), nil
}

func SaveConfig(cfg domain.ClientConfig) error {
	path, err := configFilePath("config.json")
	if err != nil {
		return err
	}

	cfg = NormalizeConfig(cfg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

func DefaultConfig() domain.ClientConfig {
	return domain.ClientConfig{
		APIURL: "https://cours-go.onrender.com",
		Theme:  "ocean",
	}
}

func NormalizeConfig(cfg domain.ClientConfig) domain.ClientConfig {
	defaults := DefaultConfig()
	if cfg.APIURL == "" {
		cfg.APIURL = defaults.APIURL
	}
	if cfg.Theme == "" {
		cfg.Theme = defaults.Theme
	}
	return cfg
}

func configFilePath(filename string) (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(baseDir, appDirectory)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, filename), nil
}
