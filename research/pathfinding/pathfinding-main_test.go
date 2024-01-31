package main

import (
	"fmt"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
	"testing"
)

func BenchmarkDijsktra4(b *testing.B) {
	m := RandomLevel(I(50), I(90), I(2000), I(2000))
	pairs := GetStartEndPairs(m, 1)
	neighbors := GetNeighborsListInt(m)

	nBigSteps := 0
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, pair := range pairs {
			_, nBigSteps, _ = findPath4(
				neighbors,
				m.PtToIndex(pair.startPt).ToInt(),
				m.PtToIndex(pair.endPt).ToInt())
		}
	}
	fmt.Println(nBigSteps)
}

func TestDijkstra1(t *testing.T) {
	m := RandomLevel(I(50), I(90), I(2000), I(2000))
	pairs := GetStartEndPairs(m, 1000)
	dijkstra1(m, pairs)
}

func TestDijkstra2(t *testing.T) {
	m := RandomLevel(I(50), I(90), I(2000), I(2000))
	pairs := GetStartEndPairs(m, 1000)
	dijkstra2(m, pairs)
}

func TestDijkstra3(t *testing.T) {
	m := RandomLevel(I(50), I(90), I(2000), I(2000))
	pairs := GetStartEndPairs(m, 10000)
	dijkstra3(m, pairs)
}

func TestDijkstra4(t *testing.T) {
	m := RandomLevel(I(50), I(90), I(2000), I(2000))
	pairs := GetStartEndPairs(m, 1000)
	dijkstra4(m, pairs)
}
