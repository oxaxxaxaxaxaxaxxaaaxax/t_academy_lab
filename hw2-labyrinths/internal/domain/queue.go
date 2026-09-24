package domain

import "errors"

var ErrQueueEmpty = errors.New("queue is empty")

type Queue struct {
	Nodes []int
}

func (q Queue) Enqueue(node int) Queue {
	q.Nodes = append(q.Nodes, node)
	return q
}

func (q Queue) Dequeue() (Queue, int, error) {
	if len(q.Nodes) == 0 {
		return Queue{}, 0, ErrQueueEmpty
	}

	c := q.Nodes[0]
	q.Nodes = q.Nodes[1:]
	return q, c, nil
}

func (q Queue) Empty() bool {
	return len(q.Nodes) == 0
}
