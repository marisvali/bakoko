package main

import (
	"fmt"
	. "playful-patterns.com/bakoko"
	. "playful-patterns.com/bakoko/ints"
)

func findPath2(m Matrix, start Pt, end Pt) (p []Pt, nBigSteps int, nSmallSteps int) {
	var queue []Pt
	//queue := make([]Pt, 0, 10000)

	//visited := make(map[Pt]bool)
	visited := Matrix{}
	visited.Init(m.NRows(), m.NCols())

	queue = append(queue, start)
	visited.Set(start.Y, start.X, I(1))
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
			if visited.Get(n.Y, n.X).Eq(I(0)) {
				queue = append(queue, n)
				parents[n] = topEl
				visited.Set(n.Y, n.X, I(1))
			}
			nSmallSteps++
		}

		// pop the first element out of the queue
		idx++
	}
	return []Pt{}, nBigSteps, nSmallSteps
}

func dijkstra2(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	// My Dijkstra 2.
	pathLens := 0
	var path []Pt

	var nBigStepsTotal, nSmallStepsTotal int
	duration := Duration(func() {
		for _, pair := range pairs {
			var nBigSteps, nSmallSteps int
			path, nBigSteps, nSmallSteps = findPath2(m, pair.startPt, pair.endPt)
			paths = append(paths, path)
			pathLens += len(path)
			nBigStepsTotal += nBigSteps
			nSmallStepsTotal += nSmallSteps
		}
	})
	fmt.Println("------------  MY DIJKSTRA 2 ------------")
	//fmt.Println("Number of paths: ", len(pairs))
	fmt.Println("Duration for n paths: ", duration)
	//fmt.Println("Duration per path: ", duration/float64(len(pairs)))
	//fmt.Println("Average path length: ", pathLens/len(pairs))
	//fmt.Println("Last path length: ", len(path))
	//fmt.Println("Number of big steps: ", nBigStepsTotal)
	fmt.Println("Average number big of steps: ", nBigStepsTotal/len(pairs))
	//fmt.Println("Number of small steps: ", nSmallStepsTotal)
	fmt.Println("Average number small of steps: ", nSmallStepsTotal/len(pairs))
	return
}
