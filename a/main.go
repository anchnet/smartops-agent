package main

import (
	"fmt"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/util/executable"
)

func main() {
	s, _ := executable.Folder()
	fmt.Println(s)

}
