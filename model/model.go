package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/steveyen/gkvlite"
)

const (
	STORE_LOCATION string = "model/persist.gkvlite"

	// Model Keys
	LAST_SUCCESS   string = "last_success"
	LAST_FAILURE   string = "last_failure"
	CURRENT_HEALTH string = "current_health"
)

var modelStore *gkvlite.Store

type SearchRecord interface {
	GetCollectionName() string
	GetLatency() time.Duration
	GetMessage() string
	GetNumberOfResults() int
}

type Record struct {
	Time            time.Time     `json:"time"`
	Latency         time.Duration `json:"latency"`
	Message         string        `json:"message"`
	NumberOfResults int           `json:"number_of_results"`
}

type Health struct {
	SuccessCount int `json:"success_count"`
	FailureCount int `json:"failure_count"`
}

// Questions
// Is the partner up or down now?
// When did the partner go down?
// When did it come back up? / How long has it been up?
// Average uptime / success rate
// Average latency time
// Average results per search

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

func persist(collectionName, key string, value []byte) (err error) {
	// Ensure collection exists before usage
	var collection *gkvlite.Collection
	if collectionExists(collectionName) {
		collection = modelStore.GetCollection(collectionName)
	} else {
		collection = modelStore.SetCollection(collectionName, nil)
	}

	err = collection.Set([]byte(key), value)
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

func SaveRecord(newRecord SearchRecord) (encodedBuiltRecord []byte, err error) {
	currentTime := time.Now()
	currentTimeStamp := currentTime.Unix()

	builtRecord := Record{
		Time:            currentTime,
		Latency:         newRecord.GetLatency(),
		Message:         newRecord.GetMessage(),
		NumberOfResults: newRecord.GetNumberOfResults(),
	}

	encodedBuiltRecord, err = json.Marshal(builtRecord)
	if err != nil {
		log.Printf("ERROR: Unable to marshal new record to json due to %v", err.Error())
		return encodedBuiltRecord, err
	}

	err = persist(newRecord.GetCollectionName(), fmt.Sprintf("%v", currentTimeStamp), encodedBuiltRecord)
	if err != nil {
		log.Printf("ERROR: Unable to persist new record due to %v", err.Error())
		return encodedBuiltRecord, err
	}

	return encodedBuiltRecord, err
}

func SaveLastSuccess(newRecord SearchRecord) (err error) {
	encodedBuiltRecord, err := SaveRecord(newRecord)
	if err != nil {
		return err
	}

	// Add new successful record
	err = persist(newRecord.GetCollectionName(), LAST_SUCCESS, encodedBuiltRecord)
	if err != nil {
		return err
	}

	return updateHealth(newRecord.GetCollectionName(), true)
}

func SaveLastFailure(newRecord SearchRecord) (err error) {
	encodedBuiltRecord, err := SaveRecord(newRecord)
	if err != nil {
		return err
	}

	// Add new failed record
	err = persist(newRecord.GetCollectionName(), LAST_FAILURE, encodedBuiltRecord)
	if err != nil {
		return err
	}

	return updateHealth(newRecord.GetCollectionName(), false)
}

func updateHealth(collectionName string, success bool) (err error) {
	retrievedHealth, err := Retrieve(collectionName, CURRENT_HEALTH)
	if err != nil {
		log.Printf("ERROR: Unable to update %v health due to %v", collectionName, err.Error())
		return err
	}

	var currentHealth Health
	if len(retrievedHealth) != 0 {
		err = json.Unmarshal(retrievedHealth, &currentHealth)
		if err != nil {
			log.Printf("ERROR: Unable to unmarshal retrieved health due to %v", err.Error())
			return err
		}
	}

	if success {
		currentHealth.SuccessCount++
	} else {
		currentHealth.FailureCount++
	}

	encodedCurrentHealth, err := json.Marshal(currentHealth)
	if err != nil {
		log.Printf("ERROR: Unable to marshal current health due to %v", err.Error())
		return err
	}

	return persist(collectionName, CURRENT_HEALTH, encodedCurrentHealth)
}

func Retrieve(collectionName, key string) (retrievedValue []byte, err error) {
	if !collectionExists(collectionName) {
		return retrievedValue, fmt.Errorf("Collection %v does not exist", collectionName)
	}

	collection := modelStore.GetCollection(collectionName)

	return collection.Get([]byte(key))
}
