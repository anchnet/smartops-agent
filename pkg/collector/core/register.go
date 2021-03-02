package core

import "github.com/anchnet/smartops-agent/pkg/collector/filter"

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
	for k, v := range catalog {
		//添加过滤
		if !filter.GetFilter().Type(k) {
			continue
		}
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
	for k, v := range pluginCatalog {
		//添加过滤
		if !filter.GetFilter().Type(k) {
			continue
		}
		pluginChecks = append(pluginChecks, v)
	}
	return pluginChecks
}
