package main

import (
	"fmt"
)

func main() {
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)
	maxNumber := 100
	for i := 1; i <= 5; i++ {
		fmt.Println("i =", maxNumber/i)
	}
}
