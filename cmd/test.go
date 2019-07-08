package main

import "fmt"

func main() {
	j, k := 1, 2
	j, k = k, j
	l := float32(k) / float32(j)
	fmt.Println("Result: ", j, k, l)
}
