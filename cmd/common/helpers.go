package common

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/config"
	"strings"
)

func SetupConfig(confFilePath string) error {
	// set the paths where a config file is expected
	if len(confFilePath) != 0 {
		// if the configuration file path was supplied on the command line,
		// add that first so it's first in line
		config.SmartOps.AddConfigPath(confFilePath)
		// If they set a config file directly, let's try to honor that
		if strings.HasSuffix(confFilePath, ".yaml") {
			config.SmartOps.SetConfigFile(confFilePath)
		}
	}
	config.SmartOps.AddConfigPath(DefaultConfPath)
	// load the configuration
	err := config.Load(config.SmartOps)
	if err != nil {
		return fmt.Errorf("unable to load SmartOps config file: %s", err)
	}
	return nil
}
