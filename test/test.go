package main

import (
	"fmt"
)

func main() {
	var x []string
	x = make([]string, 0)
	x = append(x, "1")
	x = append(x, "2")
	x = append(x, "3")

	for _, element := range x {
		fmt.Println(element)
	}
}
