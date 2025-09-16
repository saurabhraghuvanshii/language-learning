package main

import (
	"fmt"
	"time"
)

//  Add sender receiver handske
// func main() {

// 	ch := make(chan int)
// 	go func() {
// 		ch <- 43
// 	    ch <- 41
// 	}()

// 	received := <-ch
// 	// closing channel
// 	// close(ch)
// 	fmt.Println(received)

// 	// Same with buffered

// 	ch1 := make(chan int, 2)

// 	ch1 <- 5
// 	ch1 <- 8

// 	close(ch1)
// 	// syntax to cnf if channel recived or not
// 	received, ok := <- ch1

// 	if !ok {
// 		fmt.Println("channel is empty or closed")
// 	}

// 	fmt.Println(received, ok)
// }


func dummy (c chan int) {
	time.Sleep(3 * time.Second)
	<-c
}

func main() {
	ch := make(chan int)
	go dummy(ch)
	fmt.Println("waiting baby")
	ch <- 45
	fmt.Println("received")
}