package main



// a := []int{ 1,2,3,};
// b := make([] int, 5, 10 );


// how to force independent slice storage

func main() {
	a := []int{ 1,2,3,};
	b := slices.Clone(a);

	// Or

	c := make([]int, len(a))
	copy(c, a)
}
