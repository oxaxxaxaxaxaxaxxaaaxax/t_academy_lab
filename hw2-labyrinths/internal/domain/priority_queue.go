package domain

type PriorityQueue[T any] struct {
	Nodes []Node[T]
}

func (q PriorityQueue[T]) Enqueue(node Node[T]) PriorityQueue[T] {
	nodeIdx := len(q.Nodes)

	for nodeIdx > 0 && node.Priority < q.Nodes[nodeIdx-1].Priority {
		nodeIdx--
	}

	q.Nodes = append(q.Nodes, Node[T]{})
	copy(q.Nodes[nodeIdx+1:], q.Nodes[nodeIdx:len(q.Nodes)])
	q.Nodes[nodeIdx] = node
	return q
}

func (q PriorityQueue[T]) Dequeue() (PriorityQueue[T], Node[T], error) {
	if len(q.Nodes) == 0 {
		return PriorityQueue[T]{}, Node[T]{}, ErrQueueEmpty
	}

	c := q.Nodes[0]
	q.Nodes = q.Nodes[1:]
	return q, c, nil
}

func (q PriorityQueue[T]) Empty() bool {
	return len(q.Nodes) == 0
}
