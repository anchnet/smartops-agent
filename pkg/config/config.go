package config

import (
	log "github.com/cihub/seelog"
	"strings"
)

//const defaultSite = "localhost"
const defaultSite = "smartopsapi.tstack.com"

var overrideVars = map[string]interface{}{}

//SmartOps is the global configuration object
var (
	SmartOps Config
	Nginx    Config
	Mysql    Config
)

func init() {
	// Configure SmartOps global configuration
	SmartOps = NewConfig("smartops", "SO", strings.NewReplacer(".", "_"))
	// Configuration defaults
	initConfig(SmartOps)
	// Configure Nignx global configuration
	Nginx = NewConfig("nginx", "SO", strings.NewReplacer(".", "_"))
	initConfig(Nginx)
	// Configure mysql global configuration
	Mysql = NewConfig("mysql", "SO", strings.NewReplacer(".", "_"))
	initConfig(Mysql)
}

// initConfig initializes the config defaults on a config
func initConfig(config Config) {
	// Agent
	//config.SetDefault("site", defaultSite)
	config.SetDefault("site", nil)
	config.BindEnvAndSetDefault("endpoint", nil)
	config.BindEnvAndSetDefault("api_key", nil)

	// Log
	config.BindEnvAndSetDefault("log_file_max_size", "10Mb")
	config.BindEnvAndSetDefault("log_file_max_rolls", 1)
	config.BindEnvAndSetDefault("log_level", "info")
	config.BindEnvAndSetDefault("log_to_console", true)
	config.BindEnvAndSetDefault("log_format_json", false)
	// Nginx
	config.BindEnvAndSetDefault("instances.nginx_status_url", nil)
	// Mysql
	config.BindEnvAndSetDefault("instances", nil)
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
