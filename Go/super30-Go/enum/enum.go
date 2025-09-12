package main

import "fmt"

type dir int

const  (
	East dir = iota
	West
	South
	North
)

func (s dir) String() string {
	return [...]string{"East", "West", "South", "North"}[s]
}
//var x dir = 1

func main () {
	fmt.Println(dir(1))
}
