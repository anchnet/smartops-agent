package core

import (
	"fmt"
	"smart-capture/pkg/collector/check"
)

type CheckFactory func() check.Check

type GoCheckLoader struct{}

var catalog = make(map[string]CheckFactory)

func RegisterCheck(name string, c CheckFactory) {
	catalog[name] = c
}

// GetRegisteredFactoryKeys get the keys for all registered factories
func GetRegisteredFactoryKeys() []string {
	var factoryKeys []string
	for name := range catalog {
		factoryKeys = append(factoryKeys, name)
	}

	return factoryKeys
}

// Load returns a list of checks, one for every configuration instance found in `config`
func (gl *GoCheckLoader) Load(name string) (check.Check, error) {
	factory, found := catalog[name]
	if !found {
		msg := fmt.Sprintf("Check %s not found in Catalog", name)
		return nil, fmt.Errorf(msg)
	}
	return factory(), nil
}
