// Package storage implements an interface and connectors for
// reading and writing input-code pairs into a permanent storage.
package storage

import (
	"errors"

	"github.com/visheratin/url-short/config"
	"github.com/visheratin/url-short/log"
)

// instance is a private instance of Storage.
var instance Storage

// Storage is an interface that describes two main operations
// over a permanent storage:
// Store: save an input-code pair to the permanent storage;
// LoadAll: extract all input-code pairs from the permanent
// storage. This operation is performed during microservice start.
type Storage interface {
	Store(code, input string) error
	LoadAll() ([][2]string, error)
}

// initStorage reads a storage type from the config and
// creates a proper connector with required parameters.
func initStorage() error {
	config, err := config.Config()
	if err != nil {
		log.Log().Error.Println(err)
		return err
	}
	switch config.StorageType {
	case "filesystem":
		fsStorage, err := newFSStorage(config.StoragePath)
		if err != nil {
			log.Log().Error.Println(err)
			return err
		}
		instance = fsStorage
		return nil
	}
	err = errors.New("storage type was not recognized")
	log.Log().Error.Println(err)
	return err
}

// Instance returns a copy of a Storage instance.
func Instance() (Storage, error) {
	if instance == nil {
		err := initStorage()
		if err != nil {
			log.Log().Error.Println(err)
			return nil, err
		}
	}
	return instance, nil
}
