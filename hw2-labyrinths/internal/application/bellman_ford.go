package application

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"

type BellmanFordSolver struct {
}

func NewBellmanFordSolver() BellmanFordSolver {
	return BellmanFordSolver{}
}

func (b BellmanFordSolver) SolveMaze(graph domain.Graph, startCoord, goalCoord domain.Coordinates) ([]int, error) {
	edges := TransformGraphIntoEdgeList(graph)

	vertices := make([]int, graph.Width*graph.Height)
	for i := range vertices {
		vertices[i] = i
	}

	startSl := domain.CoordPair{startCoord.Row, startCoord.Col}

	start := EncodeCoordinates(startSl, graph.Width)

	parent := make([]int, graph.Width*graph.Height)
	parent[start] = -1

	adj := make([]domain.CoordPair, 0)
	for i := range edges {
		adj = append(adj, domain.CoordPair{EncodeCoordinates(edges[i].NodeDest, graph.Width), EncodeCoordinates(edges[i].NodeSource, graph.Width)})
		adj = append(adj, domain.CoordPair{EncodeCoordinates(edges[i].NodeSource, graph.Width), EncodeCoordinates(edges[i].NodeDest, graph.Width)})
	}

	distance := InitializeDistance(graph.Width * graph.Height)
	for range len(vertices) - 1 {
		for _, e := range adj {
			source := e[0]
			v := e[1]

			sourceRow := DecodeCoordinates(source, graph.Width)[0]
			sourceCol := DecodeCoordinates(source, graph.Width)[1]
			vRow := DecodeCoordinates(v, graph.Width)[0]
			vCol := DecodeCoordinates(v, graph.Width)[1]

			edgeWeight := graph.Graph[sourceRow][sourceCol].Weight +
				graph.Graph[vRow][vCol].Weight

			if distance[v] > distance[source]+edgeWeight {
				distance[v] = distance[source] + edgeWeight
				parent[v] = source
			}
		}
	}
	return parent, nil
}
