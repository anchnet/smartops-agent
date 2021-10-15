package main

import (
	"fmt"

	"github.com/anchnet/smartops-agent/pkg/util/instance/id"
)

func main() {
	s, err := id.GetInstanceId("huaweiyun")
	fmt.Println(err)
	fmt.Println(s)
}
