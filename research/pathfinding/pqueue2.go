package main

//package astar

// An item is something we manage in a priority queue.
type item2[T any] struct {
	value    T   // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
}

// A priorityQueue implements heap.Interface and holds items.
type priorityQueue2[T any] []item2[T]

func (pq priorityQueue2[T]) Len() int { return len(pq) }

func (pq priorityQueue2[T]) Less(i, j int) bool {
	// We want heap.Pop to give us the item with the highest,
	// not lowest, priority so we use greater than here.
	//return pq[i].priority.Gt(pq[j].priority)
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue2[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue2[T]) Push(x any) {
	*pq = append(*pq, x.(item2[T]))
}

func (pq *priorityQueue2[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
