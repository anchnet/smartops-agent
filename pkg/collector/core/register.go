package core

var (
	catalog       = make(map[string]Check)
	checks        = make([]Check, 0)
	pluginChecks  = make([]PluginCheck, 0)
	pluginCatalog = make(map[string]PluginCheck)
)

func RegisterCheck(c Check) {
	catalog[c.Name()] = c
}
func GetAllChecks() []Check {
	if len(checks) != 0 {
		return checks
	}
	for _, v := range catalog {
		checks = append(checks, v)
	}
	return checks
}
func RegisterPluginCheck(c PluginCheck) {
	pluginCatalog[c.PluginName()] = c
}
func GetAllPluginsCheck() []PluginCheck {
	if len(pluginChecks) != 0 {
		return pluginChecks
	}
	for _, v := range pluginCatalog {
		pluginChecks = append(pluginChecks, v)
	}
	return pluginChecks
}
