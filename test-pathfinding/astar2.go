package main

import (
	"fmt"
	"math"
	. "playful-patterns.com/bakoko"
	. "playful-patterns.com/bakoko/ints"
	"slices"
)

func heuristic(start, end, nCols int) int {
	x1 := start % nCols
	y1 := start / nCols
	x2 := end % nCols
	y2 := end / nCols
	dx := x2 - x1
	if dx < 0 {
		dx = -dx
	}
	dy := y2 - y1
	if dy < 0 {
		dy = -dy
	}
	if dx < dy {
		return dy
	} else {
		return dx
	}
}

func heuristic2(m Matrix, start, end int) int {
	return nodeDist(m.IndexToPt(I(int64(start))),
		m.IndexToPt(I(int64(end))))
}

func nodeDist(p, q Pt) int {
	d := q.Minus(p)
	return Max(d.X.Abs(), d.Y.Abs()).ToInt()
}

func getPathAstar(parents []int, end int) (path []int) {
	node := end
	for node >= 0 {
		path = append(path, node)
		node = parents[node]
	}
	slices.Reverse(path)
	return
}

func FindPath2(m Matrix, neighbors []int, start, dest int) (ppp []int, nBigSteps int, nSmallSteps int) {
	closed := make([]bool, len(neighbors)/8)
	parents := make([]int, len(neighbors)/8)
	for i := range parents {
		parents[i] = -1
	}
	cost := make([]int, len(neighbors)/8)
	for i := range cost {
		cost[i] = math.MaxInt
	}
	cost[start] = 0

	pq := &priorityQueue2[int]{}
	Init(pq, &nSmallSteps)
	Push(pq, item2[int]{value: start})

	nCols := m.NCols().ToInt()

	var idx int
	for pq.Len() > 0 {
		p := Pop(pq, &nSmallSteps).(item2[int])
		n := p.value
		if closed[n] {
			continue
		}
		if n == dest {
			// Path found
			return getPathAstar(parents, dest), idx, nSmallSteps
		}
		closed[n] = true

		idx++

		nIndex := n * 8
		ns := neighbors[nIndex : nIndex+8]
		for _, nb := range ns {
			if nb >= 0 {
				tentativeCost := cost[n] + 1
				if tentativeCost < cost[nb] {
					cost[nb] = tentativeCost
					parents[nb] = n
					Push(pq, item2[int]{
						value:    nb,
						priority: -(cost[nb] + heuristic(nb, dest, nCols)),
					})
				}
			}
		}
	}

	// No path found
	return nil, idx, nSmallSteps
}

func astarOptimized(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	neighborsInt := GetNeighborsListInt(m)

	// Internet's AStar2.
	{
		maze := graph{}
		maze.m = m
		//lastPair := pairs[len(pairs)-1]
		var path []Pt

		pathLens := 0
		var nBigStepsTotal, nSmallStepsTotal int
		duration := Duration(func() {
			for _, pair := range pairs {
				var nBigSteps, nSmallSteps int
				var pathInts []int
				pathInts, nBigSteps, nSmallSteps = FindPath2(
					m,
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
		fmt.Println("------------ ASTAR OPTIMIZED ------------")
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
