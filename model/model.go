package model

import (
	"fmt"
	"log"
	"os"

	"github.com/steveyen/gkvlite"
)

const (
	STORE_LOCATION string = "/model/persist.gkvlite"
)

var modelStore *gkvlite.Store

func init() {
	retrievedFile := findOrCreateStore()

	var err error
	modelStore, err = gkvlite.NewStore(retrievedFile)
	if err != nil {
		panic(err.Error())
	}
}

func findOrCreateStore() *os.File {
	if _, err := os.Stat(STORE_LOCATION); os.IsNotExist(err) {
		newFile, err := os.Create(STORE_LOCATION)
		if err != nil {
			panic(err.Error())
		}

		return newFile
	}

	// File already exists
	existingFile, err := os.Open(STORE_LOCATION)
	if err != nil {
		panic(err.Error())
	}

	return existingFile
}

func collectionExists(collectionName string) bool {
	existingCollections := modelStore.GetCollectionNames()
	for _, individualCollection := range existingCollections {
		if individualCollection == collectionName {
			return true
		}
	}

	return false
}

func Persist(collectionName, key, value string) (err error) {
	// Ensure collection exists before usage
	var collection *gkvlite.Collection
	if collectionExists(collectionName) {
		collection = modelStore.GetCollection(collectionName)
	} else {
		collection = modelStore.SetCollection(collectionName, nil)
	}

	err = collection.Set([]byte(key), []byte(value))
	if err != nil {
		log.Printf("ERROR: Unable to persist %v => %v due to %v", key, value, err.Error())
		return err
	}

	err = modelStore.Flush()
	if err != nil {
		log.Printf("ERROR: Unable to complete flushing collection %v due to %v", collectionName, err.Error())
		return err
	}

	return err
}

func Retrieve(collectionName, key string) (value string, err error) {
	if !collectionExists(collectionName) {
		return value, fmt.Errorf("Collection %v does not exist", collectionName)
	}

	collection := modelStore.GetCollection(collectionName)

	byteValue, err := collection.Get([]byte(key))
	if err != nil {
		return value, err
	}

	return string(byteValue), err
}
