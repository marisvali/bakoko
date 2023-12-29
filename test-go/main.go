package main

import "fmt"

type A struct {
	val int
}

func (a A) do() A {
	a.val = a.val + 1
	return a
}

func main() {
	var x A
	a := x.do().do().do()
	fmt.Println(x)
	fmt.Println(a)
}
