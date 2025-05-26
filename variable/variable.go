package main

import (
	"fmt"
)

// func printVariables()  {
// keyword identifier list type
// 	// var roomNumber, floorNumber int = 124, 4
// 	// fmt.Println(roomNumber, floorNumber)

// 	// var password string
// 	// var password = "fuck it"
// 	//short variable declaration
// 	// password := "fuck it"
// 	// fmt.Println(password)

// 	//const usinng

// 	const ver string = "123"
// 	fmt.Println(ver)
// }

func printVariables() {
	const occupancyLimit = 12

	var occupancyLimit1 uint8
	var occupancyLimit2 int64
	var occupancyLimit3 float32

	occupancyLimit1 = occupancyLimit
	occupancyLimit2 = occupancyLimit
	occupancyLimit3 = occupancyLimit

	fmt.Println(occupancyLimit1, occupancyLimit2, occupancyLimit3)
}
