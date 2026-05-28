package local

import (
	"encoding/json"
	"errors"
	"os"

	"cours-go/internal/domain"
)

func LoadState() (domain.ClientState, error) {
	path, err := configFilePath("state.json")
	if err != nil {
		return domain.ClientState{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			state := domain.ClientState{}
			return state, SaveState(state)
		}
		return domain.ClientState{}, err
	}

	var state domain.ClientState
	if err := json.Unmarshal(data, &state); err != nil {
		return domain.ClientState{}, err
	}

	return state, nil
}

func SaveState(state domain.ClientState) error {
	path, err := configFilePath("state.json")
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}
