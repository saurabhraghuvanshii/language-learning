package main

import (
	"errors"
	"fmt"
	"math"
)

func main () {
	fmt.Println("heloo baby")
	var a [5] int  // [0 0 0 0 0] //array

	a[2] = 7 // [0 0 7 0 0]

	b := [5]int{2,4,4,4,4} // this way we can put values in this [2 4 4 4 4] 

	c := []int{} // slice of int no fixed type we can add multiple value by using even append
	// slces create new array behid scene and refrence old one
	c = append(c, 33,44,44,5,5,6,)  //[33]  [33 44 44 5 5 6]
	
	// map[type value of key]types_of_value_in_the_map
	// tp create a map used built in **make** fun

	vertices := make(map[string]int)
  
	vertices["tri"] = 2
	vertices["yo"] = 1   // Output map[tri:2 yo:1]
	delete(vertices, "yo")  // Output map[tri:2]  to delete something from map


	//Loop
	for i := 0; i< 5; i ++ {
		fmt.Println(i)
	}

	// while loop using for 
	i := 0
	for i < 5 {
		fmt.Println(i)
		i++
	}
	
	//loop through array using for 
	arr := []string{"aa", "bb", "cc"}
    for index, value := range arr {  // using range we can uterate overe array
		fmt.Println("index",index,"values", value)
	}

    // loop over map
	m := make(map[string]string)
    m["a"] = "alpha"
	m["b"] = "beta"
	for key, value := range m {
		fmt.Println("key", key, "val", value )
	}

	fmt.Println(a) // this will hold  5 ints
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(vertices)
   
	result := sum(5,5)
	fmt.Println(result)

	res, err := sqrt(-16)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	// struct 

	t := person{ name: "saurabh", age: 23}
	fmt.Println(t.age, t.name)

	// Pointer 
	j := 7
	inc(&j)
	fmt.Println(j)
}


// create func

func sum (x int, y int) int  {
	return x + y
}

// In Go function can return  multiple value 

func sqrt (x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("undefined for neagtive numbers")
	}

	return math.Sqrt(x), nil
} 

// Define struct in Go
 
type person struct {
	name string
	age int
}


// Pointers in Go 

func inc(x *int) {
	*x++
}