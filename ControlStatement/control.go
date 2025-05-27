package main

// import (
// 	"fmt"
// 	"math/rand"
// )

// func main () {
//    hotelName := "The fucking hotel"
//    fmt.Println("Hotel "+ hotelName)

//    var roomsAvailable int
//    var rooms, roomOcupied int = 100 , rand.Intn(100)
//    roomsAvailable = rooms - roomOcupied
//    fmt.Println(roomsAvailable, "room available")
// }

// func main() {
// 	var agePaul, ageJohn int = rand.Intn(110), rand.Intn(110)

// 	fmt.Println("John is", ageJohn, "years old")
//     fmt.Println("Paul  is", agePaul, "years old")

// 	if ageJohn > agePaul {
// 		fmt.Println("John is older than John")
// 	} else {
// 		fmt.Println("Paul is younger than you, or both have same age")
// 	}
// }

func main() {
	var gender = "female"
   
	switch gender {
	case "male":
		println("you male")
	case "female":
		println("your female")
	default:
		println("gay bsdk ka")
	}
}
