package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Settings struct {
	API         string `json:"API"`
	Domain      string `json:"domain"`
	Subdomain   string `json:"subdomain"`
	Token       string `json:"token"`
	StartServer bool   `json:"httpServer"`
	ServerType  string `json:"ServerType"`
	Port        string `json:"httpPort"`
}

func LoadSettings(altPath string) (*Settings, error) {
	var file []byte
	var err error
	paths := []string{"/usr/local/etc/godyndns.json", "./godyndns.json"}
	if altPath != "" {
		// Use the explicitly specified path
		file, err = ioutil.ReadFile(altPath)
	} else {
		// Try all default paths
		for _, path := range paths {
			file, err = ioutil.ReadFile(path)
			if err == nil {
				break
			}
		}
	}
	if file == nil {
		return nil, errors.New("No File found")
	}

	var setting Settings
	err = json.Unmarshal(file, &setting)
	return &setting, err
}
