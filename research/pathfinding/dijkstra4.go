package main

import (
	"fmt"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
	"slices"
)

func getPath4(parents []int, end int) (path []int) {
	node := end
	for node >= 0 {
		path = append(path, node)
		node = parents[node]
	}
	slices.Reverse(path)
	return
}

func findPath4(neighbors []int, start int, end int) (p []int, nBigSteps int, nSmallSteps int) {
	// Initialize our structures.
	// Preallocation speeds things up (a little) but it needs to be done with
	// a parameter sent to the function. So I'll see about this later.
	//queue := make([]int, 100)
	var queue []int
	visited := make([]bool, len(neighbors)/8)
	parents := make([]int, len(neighbors)/8)
	for i := range parents {
		parents[i] = -1
	}

	// Process the start element.
	queue = append(queue, start)
	visited[start] = true

	idx := 0
	for idx < len(queue) {
		nBigSteps++
		// peek the first element from the queue
		topEl := queue[idx]
		if topEl == end {
			return getPath4(parents, end), nBigSteps, nSmallSteps
		}

		nIndex := topEl * 8
		ns := neighbors[nIndex : nIndex+8]
		for _, n := range ns {
			if n >= 0 && !visited[n] {
				queue = append(queue, n)
				parents[n] = topEl
				visited[n] = true
			}
			nSmallSteps++
		}

		// pop the first element out of the queue
		idx++
	}
	return []int{}, nBigSteps, nSmallSteps
}

func dijkstra4(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	neighborsInt := GetNeighborsListInt(m)

	// My Dijkstra 4.
	{
		pathLens := 0
		var path []Pt

		var nBigStepsTotal, nSmallStepsTotal int
		duration := Duration(func() {
			for _, pair := range pairs {
				var nBigSteps, nSmallSteps int
				var pathInts []int
				pathInts, nBigSteps, nSmallSteps = findPath4(
					neighborsInt,
					m.PtToIndex(pair.startPt).ToInt(),
					m.PtToIndex(pair.endPt).ToInt())

				// turn the path of ints back to path of points
				path = []Pt{}
				for i := range pathInts {
					path = append(path, m.IndexToPt(I(int64(pathInts[i]))))
				}

				paths = append(paths, path)
				pathLens += len(path)
				nBigStepsTotal += nBigSteps
				nSmallStepsTotal += nSmallSteps
			}
		})
		fmt.Println("------------  MY DIJKSTRA 4 ------------")
		//fmt.Println("Number of paths: ", len(pairs))
		fmt.Println("Duration for n paths: ", duration)
		//fmt.Println("Duration per path: ", duration/float64(len(pairs)))
		//fmt.Println("Average path length: ", pathLens/len(pairs))
		//fmt.Println("Last path length: ", len(path))
		//fmt.Println("Number of big steps: ", nBigStepsTotal)
		fmt.Println("Average number big of steps: ", nBigStepsTotal/len(pairs))
		//fmt.Println("Number of small steps: ", nSmallStepsTotal)
		fmt.Println("Average number small of steps: ", nSmallStepsTotal/len(pairs))
	}
	return
}
