// Package config contains configuration fields for the URL
// shortening microservice. It is designed as a singleton that has
// to be initialized at the beginning of the program execution.
package config

import (
	"encoding/json"
	"errors"
	"os"
)

// instance is a private instance of Configuration.
var instance Configuration

// Configuration is a container for configuration fields
// of the URL shortening microservice.
type Configuration struct {
	Init        bool
	Port        int
	StorageType string
	StoragePath string
	LogPath     string
	Prefix      string
	CodeLength  int
}

// Init reads configuration file and creates a Configuration instance.
// This method has to be called in the beginning of the execution because
// other packages rely on Configuration.
func Init(path string) error {
	var err error
	if path == "" {
		err = errors.New("configuration file path is not set")
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&instance)
	if err != nil {
		return err
	}
	(&instance).Init = true
	return nil
}

// Config returns a copy of a Configuration instance.
func Config() (Configuration, error) {
	if instance.Init {
		return instance, nil
	} else {
		err := errors.New("configuration was not initialized")
		return Configuration{}, err
	}
}
