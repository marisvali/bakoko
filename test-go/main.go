package main

import "fmt"

type A struct {
	val int
}

func do(x *[]int) {
	*x = append(*x, 13)

	//x[2] = 13
}

func main() {
	var x []int
	x = append(x, 1)
	x = append(x, 3)
	x = append(x, 5)
	x = append(x, 7)
	do(&x)
	fmt.Println(x)
}
