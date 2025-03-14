package server

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type session struct {
	path string

	DeviceCode   string `json:"device_code"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func loadSession(path string) (*session, error) {
	path = filepath.Clean(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	s.path = path
	return &s, nil
}

func (s *session) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return err
	}
	return nil
}
