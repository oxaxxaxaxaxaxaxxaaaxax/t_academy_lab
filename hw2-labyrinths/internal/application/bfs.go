package application

import (
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type BFSSolver struct{}

func NewBFSSolver() BFSSolver {
	return BFSSolver{}
}

func (BFSSolver) SolveMaze(graph domain.Graph, startCoord, goalCoord domain.Coordinates) ([]int, error) {
	edges := TransformGraphIntoEdgeList(graph)
	frontier := domain.Queue{Nodes: make([]int, 0)}
	start := EncodeCoordinates(domain.CoordPair{startCoord.Row, startCoord.Col}, graph.Width)
	parent := make([]int, graph.Width*graph.Height)
	for i := range parent {
		parent[i] = -1
	}
	frontier = frontier.Enqueue(start)

	adj := make([]domain.CoordPair, 0, len(edges)*2)
	for _, edge := range edges {
		src := EncodeCoordinates(edge.NodeSource, graph.Width)
		dst := EncodeCoordinates(edge.NodeDest, graph.Width)
		adj = append(adj, domain.CoordPair{src, dst}, domain.CoordPair{dst, src})
	}

	var vertex int
	var err error

	for !frontier.Empty() {
		frontier, vertex, err = frontier.Dequeue()
		if err != nil {
			return []int{}, err
		}
		for _, e := range adj {
			source := e[0]
			v := e[1]
			if vertex != source {
				continue
			}
			if v == start {
				continue
			}
			if parent[v] != -1 {
				continue
			}
			frontier = frontier.Enqueue(v)
			parent[v] = source
		}
	}
	return parent, nil
}
