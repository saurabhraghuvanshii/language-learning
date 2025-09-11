package main

import (
	"fmt"

	
)

func main () {
	fmt.Println(sum(5,7))
	fmt.Println(add(5,8))
	sum, sub := cal(6,7)
	fmt.Println(sum, "", sub)

	fmt.Println(cal2(6, 8))


	// Anonymous functions  First type 
	su , sb := func (a, b int ) (sum int , sub int)  {
		sum = a + b
		sub = a - b 
		return
	}(2,4)

	// Second type

	su1 , sb1 := func (a, b int ) ( int , int)  {
		sum := a + b
		sub := a - b 
		return sub, sum
	}(2,4)

	fmt.Println(su, sb)
	fmt.Println(su1, sb1)



	ans := calculator(1,3, sum1)
	fmt.Println(ans)
	ans1 := calculator(2,4, sub1)
    fmt.Println(ans1)

	double := multiplier(2)
	tripple := multiplier(3)

	fmt.Println(double(3))
	fmt.Println(tripple(4))

	divide := divider(2,3)
	fmt.Println(divide(4,5))
}

func sum ( a int  , b int ) int {
	return a + b 
}

// Variant  You can define multiple parameters of the same type in a single declaration:

func add ( a , b int )  int {
	return a + b
}

// Return Types  You can return multiple values from a function

func cal ( a, b int) ( int , int) {
	return a+b, a-b
}


// ### Named Return Values

// You can also define named return values in a function. 
// This allows you to specify the return values in the function signature, making the code easier to read and understand:

func cal2 (a, b int) ( sum , sub int) {
	sum = a +a +b ;
	sub = a - b ;
	return
}


// Function  as arguments 
// Function can take input other fun as argument.

func sum1( a, b int) int {
	return a + b
}

func sub1(a, b int ) int {
	return a - b
}

func calculator(a int, b int, fn func(int , int) int) int {
	return fn(a,b)
}

// Returning functions

func multiplier(fact int) func (int) int {
	return  func(a int) int {
		return a * fact
	}
}

func divider ( a, b float32) func ( float32, float32) float32 {
	return func(f , h float32) float32 {
		return  a + b / f + h
	}
}
