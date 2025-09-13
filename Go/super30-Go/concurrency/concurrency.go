package main

import (
	
	"fmt"
	//"time"
)

// func main () {

//    // no goroutines
//   // expensiveOp("first")
//    // expensiveOp("second")

//    // with Goroutines

//    go expensiveOp("first")
//    go expensiveOp("second")

//    time.Sleep(time.Second * 2)

//    fmt.Println("Done")
// }

// No goroutines

func expensiveOp(str string) {
	for i:= range 3 {
		fmt.Println(str, "-", i)
	}
}




// Channel 

// func main() {
//     c := make(chan int)

// 	go func() {
// 		var sum = 0
// 		for i:= range 100 {
// 			sum += i
// 		}
// 		c <- sum
// 	}()

// 	x := <-c

// 	fmt.Println(x)
// }


func main () {
	var c = make(chan bool)

	go func ()  {
		sum := 0
		for i := range 100 {
			sum +=i
		}
        fmt.Println(sum)
		c <- true
	}()

	fmt.Println("main function ended")
	<-c 
}


