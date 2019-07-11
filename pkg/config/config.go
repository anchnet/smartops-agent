package config

import (
	//"github.com/anchnet/smartops-agent/pkg/util/log"
	log "github.com/cihub/seelog"
	"strings"
)

const defaultDomain = "192.168.101.6:8100"
const defaultWsSite = "ws://" + defaultDomain + "/monitor"
const defaultSiteOri = "http://" + defaultDomain + "/"

var overrideVars = map[string]interface{}{}

//SmartCat is the global configuration object
var (
	Smartcat Config
)

func init() {
	// Configure Smartcat global configuration
	Smartcat = NewConfig("smartcat", "SC", strings.NewReplacer(".", "_"))
	// Configuration defaults
	initConfig(Smartcat)
}

// initConfig initializes the config defaults on a config
func initConfig(config Config) {
	// Agent
	config.SetDefault("domain", defaultDomain)
	config.SetDefault("ws_site", defaultWsSite)
	config.SetDefault("site_ori", defaultSiteOri)
	config.BindEnvAndSetDefault("endpoint", nil)

	// Log
	config.BindEnvAndSetDefault("log_file_max_size", "10Mb")
	config.BindEnvAndSetDefault("log_file_max_rolls", 1)
	config.BindEnvAndSetDefault("log_level", "info")
	config.BindEnvAndSetDefault("log_to_console", true)
	config.BindEnvAndSetDefault("log_format_json", false)
}

func findUnknownKeys(config Config) []string {
	var unknownKeys []string
	knownKeys := config.GetKnownKeys()
	loadedKeys := config.AllKeys()
	for _, key := range loadedKeys {
		if _, found := knownKeys[key]; !found {
			// Check if any subkey terminated with a '.*' wildcard is marked as known
			// e.g.: apm_config.* would match all sub-keys of apm_config
			splitPath := strings.Split(key, ".")
			for j := range splitPath {
				subKey := strings.Join(splitPath[:j+1], ".") + ".*"
				if _, found = knownKeys[subKey]; found {
					break
				}
			}
			if !found {
				unknownKeys = append(unknownKeys, key)
			}
		}
	}
	return unknownKeys
}

func Load(config Config) error {
	if err := config.ReadInConfig(); err != nil {
		log.Warnf("Error loading config: %v", err)
		return err
	}

	for _, key := range findUnknownKeys(config) {
		log.Warnf("Unknown key in config file: %v", key)
	}

	applyOverrides(config)
	return nil
}

// applyOverrides overrides config variables.
func applyOverrides(config Config) {
	for k, v := range overrideVars {
		config.Set(k, v)
	}
}
