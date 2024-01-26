package main

import (
	"fmt"
	. "playful-patterns.com/bakoko"
	. "playful-patterns.com/bakoko/ints"
)

func main() {
	//m := RandomLevel(I(50), I(90), I(1000), I(1000))
	//m := ManualLevel()
	m := RandomLevel(I(20), I(40), I(200), I(200))
	nPaths := 10000
	pairs := GetStartEndPairs(m, nPaths)

	paths1 := astarUnmodified(m, pairs)
	astarOptimized(m, pairs)
	paths2 := dijkstra4(m, pairs)
	paths3 := dijkstra5(m, pairs)

	if !PathsAreTheSame(m, paths1, paths2) {
		return
	}
	if !PathsAreTheSame(m, paths1, paths3) {
		return
	}
	fmt.Println()
	fmt.Println("Paths are good!")
	fmt.Println()
}
