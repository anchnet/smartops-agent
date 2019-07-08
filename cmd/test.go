package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	info, _ := cpu.Info()

	for _, i := range info {
		fmt.Println(i.Cores)
	}
}
