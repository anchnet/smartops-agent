package check

type Loader interface {
	Load(name string) (Check, error)
}
