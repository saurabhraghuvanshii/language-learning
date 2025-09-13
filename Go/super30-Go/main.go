package main

import (
	"fmt"
	// "math/rand"
)

// func main() {
// 	// println("hello world")

// 	fmt.Println(rand.Int31n(10))

// 	var a int = 42
// 	b := 5
// 	c := "saurbh"

// 	fmt.Println(a, b, c)

// 	var sum = 0
// 	for i := 0; i < 10; i++ {
// 		sum += i
// 	}
// 	fmt.Println(sum)
  
// 	var sum1 = 0
// 	for i := range 10  {
// 		// fmt.Println(i)
//         sum1 += i
// 	}
// 	fmt.Println(sum1)
// }

// func main() {
// 	var gender string

// 	fmt.Scan(&gender)
    
// 	switch gender {
// 	case "male": 
// 	   fmt.Println("ranid aadmi")
// 	case "female" :
// 		fmt.Println("auraat")
// 	}
// }


// func main () {
// 	var a[5]int = [5]int{1,2,3,4,5}
// 	var sum = 0
// 	for i := range a {
// 		sum += a[i]
// 	}
// 	fmt.Println(sum)
// }

// func main () {
// 	var arr[5]int = [5]int{1,2,3,4,5}

// 	slices := []int{1,2,3}
// 	fmt.Println(slices)

// 	slices = append(slices, 4,5,6)
// 	slices1 := append(arr[:], 5)
// 	slices2 := append(slices,6)

// 	fmt.Print(arr, slices, slices1, slices2)
// }

// map

// func main () {
// 	m := make(map[string]int)

//     m["saurabh"] = 20
// 	m["shivam"] = 21

// 	fmt.Println(m)  // Output map[saurabh:20]

// 	for key, val := range m {
// 		fmt.Println(key, val) // Output saurabh 20
// 	}

	
// }


// NExt Permutation

func nextPer(nums []int) {
	idx := -1
	n := len(nums)

	for i := n-2; i>=0; i--{
		if nums[i] < nums[i+1] {
			idx = i
			break
		}
	}

	if idx == -1 {
		reverse(nums, 0, n-1)
		return
	}

	for i := n-1; i > idx; i-- {
		if nums[i] > nums[idx] {
			nums[i], nums[idx]  = nums[idx], nums[i]
			break
		}
	}

	reverse(nums, idx+1, n-1)
	
}


func reverse(nums []int, l, r int){
	for l < r {
		nums[l], nums[r] = nums[r], nums[l]
		l++
		r--
	}
}

func main () {
	nums := []int{2,1,5,4,3,0,0,0}
	nums1 := []int{1,2,3,4,5}
	nextPer(nums)
	nextPer(nums1)
	fmt.Println(nums)
	fmt.Println(nums1)
}
