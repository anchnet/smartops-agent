package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	fmt.Println(cpu.Counts(false))
	times, _ := cpu.Times(false)
	t := times[0]
	cores, _ := cpu.Counts(false)
	fmt.Println(t.Total())
	fmt.Println(t.Total() / float64(cores))
}
