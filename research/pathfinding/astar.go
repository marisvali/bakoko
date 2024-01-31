package main

import (
	"fmt"
	"github.com/fzipp/astar"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
)

// The Graph interface is the minimal interface a graph data structure
// must satisfy to be suitable for the A* algorithm.
type Graph[Node any] interface {
	// Neighbours returns the neighbour nodes of node n in the graph.
	Neighbours(n Node) []Node
}

// A CostFunc is a function that returns a cost for the transition
// from node a to node b.
type CostFunc[Node any] func(a, b Node) float64

// A Path is a sequence of nodes in a graph.
type Path[Node any] []Node

// newPath creates a new path with one start node. More nodes can be
// added with append().
func newPath[Node any](start Node) Path[Node] {
	return []Node{start}
}

// last returns the last node of path p. It is not removed from the path.
func (p Path[Node]) last() Node {
	return p[len(p)-1]
}

// cont creates a new path, which is a continuation of path p with the
// additional node n.
func (p Path[Node]) cont(n Node, nSteps *int) Path[Node] {
	newPath := make([]Node, len(p), len(p)+1)
	copy(newPath, p)
	newPath = append(newPath, n)
	*nSteps += len(newPath)
	return newPath
}

// Cost calculates the total cost of path p by applying the cost function d
// to all path segments and returning the sum.
func (p Path[Node]) Cost(d CostFunc[Node]) (c float64) {
	for i := 1; i < len(p); i++ {
		//c.Add(d(p[i-1], p[i]))
		c += d(p[i-1], p[i])
	}
	return c
}

// FindPath finds the shortest path between start and dest in graph g
// using the cost function d and the cost heuristic function h.
func FindPath[Node comparable](g Graph[Node], start, dest Node, d, h CostFunc[Node]) (ppp Path[Node], nBigSteps int, nSmallSteps int) {
	closed := make(map[Node]bool)

	pq := &priorityQueue[Path[Node]]{}
	Init(pq, &nSmallSteps)
	Push(pq, &item[Path[Node]]{value: newPath(start)})

	var idx int
	for pq.Len() > 0 {
		p := Pop(pq, &nSmallSteps).(*item[Path[Node]]).value
		n := p.last()
		if closed[n] {
			continue
		}
		if n == dest {
			// Path found
			return p, idx, nSmallSteps
		}
		closed[n] = true

		idx++

		for _, nb := range g.Neighbours(n) {
			newPath := p.cont(nb, &nSmallSteps)
			Push(pq, &item[Path[Node]]{
				value: newPath,
				//priority: newPath.Cost(d).Plus(h(nb, dest)).Negative(),
				priority: -(newPath.Cost(d) + (h(nb, dest))),
			})
			nSmallSteps += len(newPath)
		}
	}

	// No path found
	return nil, idx, nSmallSteps
}

// nodeDist is our cost function. We use points as nodes, so we
// calculate their Euclidean distance.
func nodeDistPt(p, q Pt) float64 {
	d := q.Minus(p)
	// Yes, Euclidian distance is what we aim for.
	// But since we're just comparing the distances, do we really need to sqrt?
	//return d.X.Sqr().Plus(d.Y.Sqr()).Sqrt()
	//return d.X.Sqr().Plus(d.Y.Sqr())

	// I'm getting errors, let's go back to float64 distance
	//return math.Sqrt(d.X.Sqr().ToFloat64() + d.Y.Sqr().ToFloat64())

	// My Dijkstra implementation implicitly assumes that the distance
	// between two neighbors is 1, even if they are on the diagonal.
	// Maybe this is the reason why AStar and Dijkstra are giving me different
	// results.
	// Make AStar's dist function behave the same.
	return Max(d.X.Abs(), d.Y.Abs()).ToFloat64()
	// Yes, this was the problem.
}

type graph struct {
	m Matrix
}

// Neighbours implements the astar.Graph[Node] interface (with Node = image.Point).
func (g graph) Neighbours(p Pt) []Pt {
	offsets := []Pt{
		{I(0), I(-1)}, // North
		{I(1), I(0)},
		{I(0), I(1)},
		{I(-1), I(0)},
		//diagonals
		{I(-1), I(-1)}, // North
		{I(1), I(1)},
		{I(-1), I(1)},
		{I(1), I(-1)},
	}
	res := make([]Pt, 0, 8)
	for _, off := range offsets {
		q := p.Plus(off)
		if g.isFreeAt(q) {
			res = append(res, q)
		}
		//else {
		//	if g.m.Get(q.Y, q.X).Neq(I(0)) {
		//		fmt.Println("AStar found obstacle")
		//	}
		//}
	}
	return res
}

func (g graph) isFreeAt(p Pt) bool {
	return g.isInBounds(p) &&
		g.m.Get(p.Y, p.X).Eq(I(0))
}

func (g graph) isInBounds(p Pt) bool {
	return p.Y.Geq(I(0)) && p.X.Geq(I(0)) &&
		p.Y.Lt(g.m.NRows()) &&
		p.X.Lt(g.m.NCols())
}

func astarWithStepsCounted(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	// Internet's AStar.
	maze := graph{}
	maze.m = m
	//lastPair := pairs[len(pairs)-1]
	var path Path[Pt]

	pathLens := 0
	var nBigStepsTotal, nSmallStepsTotal int
	duration := Duration(func() {
		for _, pair := range pairs {
			var nBigSteps, nSmallSteps int
			path, nBigSteps, nSmallSteps = FindPath[Pt](maze, pair.startPt, pair.endPt, nodeDistPt, nodeDistPt)
			paths = append(paths, path)
			pathLens += len(path)
			nBigStepsTotal += nBigSteps
			nSmallStepsTotal += nSmallSteps
		}
	})
	fmt.Println("------------ INTERNET ASTAR ------------")
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

func astarUnmodified(m Matrix, pairs []StartEnd) (paths [][]Pt) {
	//Internet's AStar, unmodified.
	maze := graph{}
	maze.m = m
	//lastPair := pairs[len(pairs)-1]
	var path astar.Path[Pt]

	pathLens := 0
	duration := Duration(func() {
		for _, pair := range pairs {

			path = astar.FindPath[Pt](maze, pair.startPt, pair.endPt, nodeDistPt, nodeDistPt)
			paths = append(paths, path)
			pathLens += len(path)
		}
	})
	fmt.Println("------------ INTERNET ASTAR UMODIFIED ------------")
	fmt.Println("Number of dpaths: ", len(pairs))
	fmt.Println("Duration for n dpaths: ", duration)
	fmt.Println("Duration per path: ", duration/float64(len(pairs)))
	fmt.Println("Average path length: ", pathLens/len(pairs))
	fmt.Println("Last path length: ", len(path))
	return
}
