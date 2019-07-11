package core

import (
	"gitlab.51idc.com/smartops/smartops-agent/pkg/collector/check"
)

type GoCheckLoader struct{}

var catalog = make(map[string]check.Check)

func RegisterCheck(name string, c check.Check) {
	catalog[name] = c
}

// Load returns a list of checks, one for every configuration instance found in `config`
func LoadChecks() []check.Check {
	var checks []check.Check
	for _, v := range catalog {
		checks = append(checks, v)
	}
	return checks
}
