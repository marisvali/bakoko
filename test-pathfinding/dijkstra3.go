package main

import (
	"fmt"
	. "playful-patterns.com/bakoko"
	. "playful-patterns.com/bakoko/ints"
	"slices"
)

func getPath3(parents []Int, end Int) (path []Int) {
	node := end
	for node.IsNonNegative() {
		path = append(path, node)
		node = parents[node.ToInt()]
	}
	slices.Reverse(path)
	return
}

func findPath3(neighbors []Int, start Int, end Int) (p []Int, nBigSteps int, nSmallSteps int) {
	// Initialize our structures.
	var queue []Int
	visited := make([]bool, len(neighbors)/8)
	parents := make([]Int, len(neighbors)/8)
	for i := range parents {
		parents[i] = I(-1)
	}

	// Process the start element.
	queue = append(queue, start)
	visited[start.ToInt()] = true

	idx := 0
	for idx < len(queue) {
		nBigSteps++
		// peek the first element from the queue
		topEl := queue[idx]
		if topEl.Eq(end) {
			return getPath3(parents, end), nBigSteps, nSmallSteps
		}

		nIndex := topEl.ToInt() * 8
		ns := neighbors[nIndex : nIndex+8]
		for _, n := range ns {
			if n.IsNonNegative() && !visited[n.ToInt()] {
				queue = append(queue, n)
				parents[n.ToInt()] = topEl
				visited[n.ToInt()] = true
			}
			nSmallSteps++
		}

		// pop the first element out of the queue
		idx++
	}
	return []Int{}, nBigSteps, nSmallSteps
}

func dijkstra3(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	neighbors := GetNeighborsList(m)

	// My Dijkstra 3.
	{
		pathLens := 0
		var path []Pt

		var nBigStepsTotal, nSmallStepsTotal int
		duration := Duration(func() {
			for _, pair := range pairs {
				var nBigSteps, nSmallSteps int
				var pathInts []Int
				pathInts, nBigSteps, nSmallSteps = findPath3(
					neighbors,
					m.PtToIndex(pair.startPt),
					m.PtToIndex(pair.endPt))

				// turn the path of ints back to path of points
				path = []Pt{}
				for i := range pathInts {
					path = append(path, m.IndexToPt(pathInts[i]))
				}

				paths = append(paths, path)
				pathLens += len(path)
				nBigStepsTotal += nBigSteps
				nSmallStepsTotal += nSmallSteps
			}
		})
		fmt.Println("------------  MY DIJKSTRA 3 ------------")
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
