package common

import (
	"fmt"
	"gitlab.51idc.com/smartops/smartops-agent/pkg/config"
	"strings"
)

func SetupConfig(confFilePath string) error {
	// set the paths where a config file is expected
	if len(confFilePath) != 0 {
		// if the configuration file path was supplied on the command line,
		// add that first so it's first in line
		config.Smartcat.AddConfigPath(confFilePath)
		// If they set a config file directly, let's try to honor that
		if strings.HasSuffix(confFilePath, ".yaml") {
			config.Smartcat.SetConfigFile(confFilePath)
		}
	}
	config.Smartcat.AddConfigPath(DefaultConfPath)
	// load the configuration
	err := config.Load(config.Smartcat)
	if err != nil {
		return fmt.Errorf("unable to load Smartcat config file: %s", err)
	}
	return nil
}
