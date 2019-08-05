package core

var (
	catalog = make(map[string]Check)
)

func RegisterCheck(name string, c Check) {
	catalog[name] = c
}
func GetAllChecks() []Check {
	var checks []Check
	for _, v := range catalog {
		checks = append(checks, v)
	}
	return checks
}
