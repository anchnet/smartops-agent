package core

import (
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"
)

type GoCheckLoader struct{}

var catalog = make(map[string]check.Check)

func RegisterCheck(name string, c check.Check) {
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
func LoadChecks() []check.Check {
	var checks []check.Check
	for _, v := range catalog {
		checks = append(checks, v)
	}
	return checks
}
