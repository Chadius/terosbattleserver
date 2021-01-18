package terosbattleserver

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// SquaddieRepository will interact with external devices to manage Squaddies.
type SquaddieRepository struct {
	squaddiesByName map[string]Squaddie
}

// NewSquaddieRepository generates a pointer to a new Squaddie.
func NewSquaddieRepository() *SquaddieRepository {
	repository := SquaddieRepository{
		map[string]Squaddie{},
	}
	return &repository
}

type unmarshalFunc func([]byte, interface{}) error

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *SquaddieRepository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *SquaddieRepository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *SquaddieRepository) addSource(data []byte, unmarshal unmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfSquaddies []Squaddie
	unmarshalError = unmarshal(data, &listOfSquaddies)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for _, squaddieToAdd := range listOfSquaddies {
		squaddieErr := CheckSquaddieForErrors(&squaddieToAdd)
		if squaddieErr != nil {
			return false, squaddieErr
		}
		squaddieToAdd.SetHPToMax()
		repository.squaddiesByName[squaddieToAdd.Name] = squaddieToAdd
	}

	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *SquaddieRepository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByName)
}

// GetByName retrieves a Squaddie by name
func (repository *SquaddieRepository) GetByName(squaddieName string) *Squaddie {
	squaddie, squaddieExists := repository.squaddiesByName[squaddieName]
	if !squaddieExists {
		return nil
	}
	return &squaddie
}

// MarshalSquaddieIntoJSON converts the given Squaddie into JSON.
func (repository *SquaddieRepository) MarshalSquaddieIntoJSON(squaddie *Squaddie) ([]byte, error) {
	return json.Marshal(squaddie)
}

// AddInnatePowersToSquaddie will find the Squaddie's powers from the PowerRepository and give them to it.
//   Returns the number of powers added.
//   Throws an error if it can't find a power.
func (repository *SquaddieRepository) AddInnatePowersToSquaddie(squaddie *Squaddie, powerRepo *PowerRepository) (int, error) {
	numberOfPowersAdded := 0

	for _, powerIDName := range squaddie.PowerIDNames {
		powerToAdd := powerRepo.GetByName(powerIDName.Name)
		if powerToAdd == nil {
			return numberOfPowersAdded, fmt.Errorf("squaddie '%s' tried to add Power '%s' but it does not exist", squaddie.Name, powerIDName.Name)
		}
		squaddie.GainInnatePower(powerToAdd)
		numberOfPowersAdded = numberOfPowersAdded + 1
	}

	return numberOfPowersAdded, nil
}
