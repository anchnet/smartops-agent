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
		return fmt.Errorf("unable to load CloudOps config file: %s", err)
	}
	return nil
}

func SetupNgxConfig(confFilePath string) error {
	config.Nginx.AddConfigPath(confFilePath)
	if err := config.Load(config.Nginx); err != nil {
		return fmt.Errorf("unable to load Nginx config file %s:", err)
	}
	return nil
}

func SetUpMysqlConfig(confFilePath string) error {
	config.Mysql.AddConfigPath(confFilePath)
	if err := config.Load(config.Mysql); err != nil {
		return fmt.Errorf("unable to load Mysql config file %s:", err)
	}
	return nil
}
