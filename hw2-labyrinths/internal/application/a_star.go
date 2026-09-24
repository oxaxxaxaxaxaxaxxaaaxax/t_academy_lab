package application

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"

type AStarSolver struct{}

func NewAStarSolver() AStarSolver {
	return AStarSolver{}
}

func (AStarSolver) SolveMaze(graph domain.Graph, startCoord, goalCoord domain.Coordinates) ([]int, error) {
	start := EncodeCoordinates(domain.CoordPair{startCoord.Row, startCoord.Col}, graph.Width)
	goal := EncodeCoordinates(domain.CoordPair{goalCoord.Row, goalCoord.Col}, graph.Width)
	edges := TransformGraphIntoEdgeList(graph)

	var frontier domain.PriorityQueue[int]

	movementCost := make([]int, graph.Width*graph.Height)
	for i := range movementCost {
		movementCost[i] = infiniteWeight
	}
	movementCost[start] = 0

	parent := make([]int, graph.Width*graph.Height)
	for i := range parent {
		parent[i] = -1
	}

	//in this realization Node.Value is node number
	frontier = frontier.Enqueue(domain.Node[int]{Priority: 0, Value: start})

	adj := make([]domain.CoordPair, 0, len(edges)*2)
	for _, edge := range edges {
		src := EncodeCoordinates(edge.NodeSource, graph.Width)
		dst := EncodeCoordinates(edge.NodeDest, graph.Width)
		adj = append(adj, domain.CoordPair{src, dst}, domain.CoordPair{dst, src})
	}

	var vertex domain.Node[int]
	var err error

	for !frontier.Empty() {
		frontier, vertex, err = frontier.Dequeue()
		if err != nil {
			return []int{}, err
		}
		if vertex.Value == goal {
			break
		}
		for _, e := range adj {
			source := e[0]
			v := e[1]
			if vertex.Value != source {
				continue
			}
			if v == start {
				continue
			}
			sourceCoord := DecodeCoordinates(source, graph.Width)
			vCoord := DecodeCoordinates(v, graph.Width)

			sourceRow, sourceCol := sourceCoord[0], sourceCoord[1]
			vRow, vCol := vCoord[0], vCoord[1]

			edgeWeight := graph.Graph[sourceRow][sourceCol].Weight +
				graph.Graph[vRow][vCol].Weight

			newCost := movementCost[source] + edgeWeight
			if newCost < movementCost[v] {
				movementCost[v] = newCost
				priority := newCost + heuristic(goal, v, graph.Width)
				frontier = frontier.Enqueue(domain.Node[int]{Priority: priority, Value: v})
				parent[v] = source
			}
		}
	}
	return parent, nil
}

func heuristic(goal, currentVertex, width int) int {
	cordGoal := DecodeCoordinates(goal, width)
	cordVertex := DecodeCoordinates(currentVertex, width)
	return absInt(cordGoal[0]-cordVertex[0]) + absInt(cordGoal[1]-cordVertex[1])
}

func absInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
