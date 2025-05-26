package main

import (
	"fmt"
)

func main() {
	// now := time.Now()
	// fmt.Println(now)
	// fmt.Println("hello world ")
	// n := 2548
	// fmt.Printf("%x", n)
	// fmt.Printf("%d",n)
	b := "hello"
	for i := 0; i < len(b); i++ {
		fmt.Printf("%c\n", b[i])
	}

	printVariables()
}

