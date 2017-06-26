package main

import "errors"
import "encoding/json"
import "io/ioutil"

type settings struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	Token     string `json:"token"`
}

func loadSettings(altPath string) (*settings, error) {
	var file []byte
	var err error
	paths := []string{"/usr/local/etc/do-dyndns.json", "./do-dyndns.json"}
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

	var setting settings
	err = json.Unmarshal(file, &setting)
	return &setting, err
}
