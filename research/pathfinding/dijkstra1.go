package main

import (
	"fmt"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
	"slices"
)

func validNeighbor(m Matrix, pt Pt) bool {
	return pt.X.Geq(I(0)) &&
		pt.Y.Geq(I(0)) &&
		pt.X.Lt(m.NCols()) &&
		pt.Y.Lt(m.NRows()) &&
		m.Get(pt.Y, pt.X).Eq(I(0))
}

func getNeighbours(m Matrix, pt Pt) (ns []Pt) {
	potentials := []Pt{
		{pt.X.Minus(I(1)), pt.Y},
		{pt.X, pt.Y.Minus(I(1))},
		{pt.X.Plus(I(1)), pt.Y},
		{pt.X, pt.Y.Plus(I(1))},
		//diagonal neighbours
		{pt.X.Minus(I(1)), pt.Y.Minus(I(1))},
		{pt.X.Plus(I(1)), pt.Y.Plus(I(1))},
		{pt.X.Plus(I(1)), pt.Y.Minus(I(1))},
		{pt.X.Minus(I(1)), pt.Y.Plus(I(1))},
	}

	ns = make([]Pt, 0, 8)
	for _, p := range potentials {
		if validNeighbor(m, p) {
			ns = append(ns, p)
		}
	}
	return
}

func getPath(parents map[Pt]Pt, end Pt) (path []Pt) {
	node := end
	var ok = true
	for ok {
		path = append(path, node)
		node, ok = parents[node]
	}
	slices.Reverse(path)
	return
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

func Contains[S ~[]E, E comparable](s S, v E) (bool, int) {
	i := Index(s, v)
	if i >= 0 {
		return true, i
	} else {
		return false, len(s)
	}
}

func findPath(m Matrix, start Pt, end Pt) (p []Pt, nBigSteps int, nSmallSteps int) {
	var queue []Pt
	queue = append(queue, start)
	parents := make(map[Pt]Pt)
	idx := 0
	for idx < len(queue) {
		nBigSteps++
		// peek the first element from the queue
		topEl := queue[idx]
		if topEl.Eq(end) {
			return getPath(parents, end), nBigSteps, nSmallSteps
		}

		ns := getNeighbours(m, topEl)
		for _, n := range ns {
			contains, nStepsForContains := Contains(queue, n)
			if !contains {
				parents[n] = topEl
				queue = append(queue, n)
			}
			nSmallSteps += nStepsForContains
		}

		// pop the first element out of the queue
		idx++
	}
	return []Pt{}, nBigSteps, nSmallSteps
}

func dijkstra1(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	// My Dijkstra.
	pathLens := 0
	var path []Pt
	var nBigStepsTotal, nSmallStepsTotal int
	duration := Duration(func() {
		for _, pair := range pairs {
			var nBigSteps, nSmallSteps int
			path, nBigSteps, nSmallSteps = findPath(m, pair.startPt, pair.endPt)
			paths = append(paths, path)
			pathLens += len(path)
			nBigStepsTotal += nBigSteps
			nSmallStepsTotal += nSmallSteps
		}
	})
	fmt.Println("------------   MY DIJKSTRA  ------------")
	fmt.Println("Number of dpaths: ", len(pairs))
	fmt.Println("Duration for n dpaths: ", duration)
	fmt.Println("Duration per path: ", duration/float64(len(pairs)))
	fmt.Println("Average path length: ", pathLens/len(pairs))
	fmt.Println("Last path length: ", len(path))
	fmt.Println("Number of big steps: ", nBigStepsTotal)
	fmt.Println("Average number big of steps: ", nBigStepsTotal/len(pairs))
	fmt.Println("Number of small steps: ", nSmallStepsTotal)
	fmt.Println("Average number small of steps: ", nSmallStepsTotal/len(pairs))
	return
}
