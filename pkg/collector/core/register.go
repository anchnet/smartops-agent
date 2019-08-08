package core

var (
	catalog = make(map[string]Check)
	checks  = make([]Check, 0)
)

func RegisterCheck(name string, c Check) {
	catalog[name] = c
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
