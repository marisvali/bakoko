package main

import (
	"fmt"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
	"slices"
)

type PathfindingTest struct {
	neighbors []int
	visited   []bool
	parents   []int
	queue     []int
	nDirs     int
	m         Matrix
}

func (p *PathfindingTest) Initialize(m Matrix) {
	// Keep reference to Matrix in order to transform Pts to ints and ints to
	// Pts in the FindPath method.
	p.m = m

	// Turn matrix into an array of ints.
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
	p.nDirs = len(dirs)

	// At neighbors[i] we will find the 8 neighbors of node with index i.
	// Each neighbor is another index. If the index is -1, the neighbor is
	// invalid.
	p.neighbors = make([]int, m.NRows().Times(m.NCols()).ToInt()*len(dirs))
	for y := I(0); y.Lt(m.NRows()); y.Inc() {
		for x := I(0); x.Lt(m.NCols()); x.Inc() {
			pt := Pt{x, y}
			index := m.PtToIndex(pt).ToInt() * p.nDirs
			ns := p.neighbors[index : index+p.nDirs]
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

	// This slice should never be re-allocated.
	p.queue = make([]int, 0, m.NRows().Times(m.NCols()).ToInt())
	// These slices will never be resized.
	p.visited = make([]bool, len(p.neighbors)/p.nDirs)
	p.parents = make([]int, len(p.neighbors)/p.nDirs)
}

func (p *PathfindingTest) computePath(parents []int, end int) (path []Pt) {
	node := end
	for node >= 0 {
		path = append(path, p.m.IndexToPt(I(node)))
		node = parents[node]
	}
	slices.Reverse(path)
	return
}

func (p *PathfindingTest) FindPath(startPt, endPt Pt) []Pt {
	// Convert Pts to ints.
	start := p.m.PtToIndex(startPt).ToInt()
	end := p.m.PtToIndex(endPt).ToInt()

	// Initialize our structures.
	p.queue = p.queue[:0] // Make len(p.queue) == 0 without re-allocating.
	for i := range p.parents {
		p.parents[i] = -1
		p.visited[i] = false
	}

	// Process the start element.
	p.queue = append(p.queue, start)
	p.visited[start] = true

	idx := 0
	for idx < len(p.queue) {
		// peek the first element from the queue
		topEl := p.queue[idx]
		if topEl == end {
			return p.computePath(p.parents, end)
		}

		nIndex := topEl * p.nDirs
		ns := p.neighbors[nIndex : nIndex+p.nDirs]
		for _, n := range ns {
			if n >= 0 && !p.visited[n] {
				p.queue = append(p.queue, n)
				p.parents[n] = topEl
				p.visited[n] = true
			}
		}

		// pop the first element out of the queue
		idx++
	}
	return []Pt{}
}

func dijkstra5(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	// My Dijkstra 5.
	pathLens := 0
	var path []Pt

	var pathfinding PathfindingTest
	pathfinding.Initialize(m)

	duration := Duration(func() {
		for _, pair := range pairs {
			path = pathfinding.FindPath(pair.startPt, pair.endPt)
			paths = append(paths, path)
			pathLens += len(path)
		}
	})
	fmt.Println("------------  MY DIJKSTRA 5 ------------")
	//fmt.Println("Number of paths: ", len(pairs))
	fmt.Println("Duration for n paths: ", duration)
	//fmt.Println("Duration per path: ", duration/float64(len(pairs)))
	fmt.Println("Average path length: ", pathLens/len(pairs))
	//fmt.Println("Last path length: ", len(path))
	return
}
