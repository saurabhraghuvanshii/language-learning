package main

import (
	"fmt"
	"math/rand"
)

func main () {
   hotelName := "The fucking hotel"
   fmt.Println("Hotel "+ hotelName)

   var roomsAvailable int
   var rooms, roomOcupied int = 100 , rand.Intn(100)
   roomsAvailable = rooms - roomOcupied
   fmt.Println(roomsAvailable, "room available")
}
