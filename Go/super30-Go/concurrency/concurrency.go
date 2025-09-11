package main

import (
	"fmt"
	"time"
)

func main () {

   // no goroutines
  // expensiveOp("first")
   // expensiveOp("second")

   // with Goroutines

   go expensiveOp("first")
   go expensiveOp("second")

   time.Sleep(time.Second * 2)

   fmt.Println("Done")
}

// No goroutines

func expensiveOp(str string) {
	for i:= range 3 {
		fmt.Println(str, "-", i)
	}
}

