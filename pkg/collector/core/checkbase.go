package core

type CheckBase struct {
	checkName string
}

func NewCheckBase(checkName string) CheckBase {
	return CheckBase{checkName: checkName}
}

func (c CheckBase) String() string {
	return c.checkName
}

//
