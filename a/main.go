package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	//s, _ := executable.Folder()
	//fmt.Println(s)

	//exPath := filepath.Dir("/Users/brick/go/src/gitlab.51idc.com/smartops/smartops-agent/a/main.go")
	//exPath := filepath.Dir("a/main.go")
	//err,exPath := filepath.Abs("a/main.go")
	fmt.Println(filepath.Abs("a/main.go"))

	//fmt.Println(exPath)

}
