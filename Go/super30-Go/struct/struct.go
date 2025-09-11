package main

import ( 
	"fmt"
)

type rect struct {
	height int
	weight int 
}

type circle struct {
	radius float32
}

func main () {
	r := rect{ weight: 10, height: 14}

	fmt.Println(area(r))

	fmt.Println(peri(r))

	rad := circle{ radius: 10}
	fmt.Println(circlearea(rad))

}

func area(r rect) int {
	return r.height * r.weight
}

func peri(a rect) int {
	return 2 * (a.height + a.weight)
}

func circlearea( rad circle) float32 {
	return 3.14 * rad.radius * rad.radius
}

