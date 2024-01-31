package main

import (
	"fmt"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
)

type StartEnd struct {
	startPt Pt
	endPt   Pt
}

func GetStartEndPairs(m Matrix, nPairs int) (pairs []StartEnd) {
	for len(pairs) < nPairs {
		var pair StartEnd
		pair.startPt.X = RInt(I(0), m.NCols().Minus(I(1)))
		pair.startPt.Y = RInt(I(0), m.NRows().Minus(I(1)))
		pair.endPt.X = RInt(I(0), m.NCols().Minus(I(1)))
		pair.endPt.Y = RInt(I(0), m.NRows().Minus(I(1)))
		if m.Get(pair.startPt.Y, pair.startPt.X).IsZero() &&
			m.Get(pair.endPt.Y, pair.endPt.X).IsZero() {
			pairs = append(pairs, pair)
		}
	}
	return
}

func GetNeighborsList(m Matrix) []Int {
	// Turn matrix into an array of Ints.
	dirs := []Pt{
		// left/right
		{I(1).Negative(), I(0)},
		{I(1), I(0)},
		// top
		{I(1).Negative(), I(1).Negative()},
		{I(0), I(1).Negative()},
		{I(1), I(1).Negative()},
		// bottom
		{I(1).Negative(), I(1)},
		{I(0), I(1)},
		{I(1), I(1)},
	}

	// At neighbors[i] we will find the 8 neighbors of node with index i.
	// Each neighbor is another index. If the index is -1, the neighbor is
	// invalid.
	neighbors := make([]Int, m.NRows().Times(m.NCols()).ToInt64()*8)
	for y := I(0); y.Lt(m.NRows()); y.Inc() {
		for x := I(0); x.Lt(m.NCols()); x.Inc() {
			pt := Pt{x, y}
			index := m.PtToIndex(pt).Times(I(8))
			ns := neighbors[index.ToInt():index.Plus(I(8)).ToInt()]
			for i := range dirs {
				neighbor := pt.Plus(dirs[i])
				if m.InBounds(neighbor) && m.Get(neighbor.Y, neighbor.X).Eq(I(0)) {
					ns[i] = m.PtToIndex(neighbor)
				} else {
					ns[i] = I(-1)
				}
			}
		}
	}
	return neighbors
}

func GetNeighborsListInt(m Matrix) []int {
	// Turn matrix into an array of Ints.
	dirs := []Pt{
		// left/right
		{I(1).Negative(), I(0)},
		{I(1), I(0)},
		// top
		{I(1).Negative(), I(1).Negative()},
		{I(0), I(1).Negative()},
		{I(1), I(1).Negative()},
		// bottom
		{I(1).Negative(), I(1)},
		{I(0), I(1)},
		{I(1), I(1)},
	}

	// At neighbors[i] we will find the 8 neighbors of node with index i.
	// Each neighbor is another index. If the index is -1, the neighbor is
	// invalid.

	neighborsInt := make([]int, m.NRows().Times(m.NCols()).ToInt64()*8)
	for y := I(0); y.Lt(m.NRows()); y.Inc() {
		for x := I(0); x.Lt(m.NCols()); x.Inc() {
			pt := Pt{x, y}
			index := m.PtToIndex(pt).Times(I(8))
			ns := neighborsInt[index.ToInt():index.Plus(I(8)).ToInt()]
			for i := range dirs {
				neighbor := pt.Plus(dirs[i])
				if m.InBounds(neighbor) && m.Get(neighbor.Y, neighbor.X).Eq(I(0)) {
					ns[i] = m.PtToIndex(neighbor).ToInt()
				} else {
					ns[i] = -1
				}
			}
		}
	}
	return neighborsInt
}

func PathsAreTheSame(m Matrix, paths1 [][]Pt, paths2 [][]Pt) bool {
	// Compare paths1 from AStar and Dijkstra.
	for i := range paths1 {
		p1 := paths1[i]
		p2 := paths2[i]

		// Can't check if they're the same because in a matrix we have
		// multiple optimal paths.
		//for j := range p1 {
		//	if !p1[j].Eq(a[j]) {
		//		fmt.Println("Found error!")
		//		return
		//	}
		//}
		// Paths are just as good if they have the same length.
		// The only thing left to check is if they are correct. This means
		// consecutive nodes are adjacent and there's no obstacle at any node
		// position.
		if !goodPath(p1, m) {
			fmt.Println("Problem with Dijkstra.")
			return false
		}

		if !goodPath(p2, m) {
			fmt.Println("Problem with AStar2.")
			return false
		}

		if len(p1) != len(p2) {
			fmt.Println("Found different lengths for paths!")
			return false
		}
	}
	return true
}

func goodPath(p []Pt, m Matrix) bool {
	for i := 0; i < len(p)-1; i++ {
		p1 := p[i]
		if m.Get(p1.Y, p1.X).Neq(I(0)) {
			fmt.Println("Path goes over obstacle.")
			return false
		}

		p2 := p[i+1]
		d := p2.Minus(p1)
		adjacent := d.X.Abs().Leq(I(1)) && d.Y.Abs().Leq(I(1))
		if !adjacent {
			fmt.Println("Path contains non-adjacent consecutive points.")
			return false
		}
	}
	return true
}
