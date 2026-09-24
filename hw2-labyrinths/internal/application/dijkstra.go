package application

import (
	"math"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

const infiniteWeight = math.MaxInt32

type DijkstraSolver struct {
}

func NewDijkstraSolver() DijkstraSolver {
	return DijkstraSolver{}
}

func TransformGraphIntoEdgeList(graph domain.Graph) []domain.Edge {
	mazeGraph := graph.Graph
	edges := make([]domain.Edge, 0)

	for _, c := range mazeGraph {
		for _, cell := range c {
			edges = AppendEdge(cell, edges, graph)
		}
	}

	return edges
}

func AppendEdge(cell domain.Cell, edges []domain.Edge, graph domain.Graph) []domain.Edge {
	neighbors := getNeighborsInMaze(cell)
	for _, neighbor := range neighbors {
		edgeWeight := graph.Graph[neighbor[0]][neighbor[1]].Weight + cell.Weight
		currEdge := domain.NewEdge(cell.Row, cell.Col, neighbor[0],
			neighbor[1], edgeWeight)
		if len(edges) == 0 {
			edges = append(edges, currEdge)
		}
		edges = AppendIfUnique(currEdge, edges)
	}
	return edges
}

func AppendIfUnique(currEdge domain.Edge, edges []domain.Edge) []domain.Edge {
	for _, edge := range edges {
		if currEdge.IsEqual(edge) {
			return edges
		}
	}
	edges = append(edges, currEdge)
	return edges
}

func getNeighborsInMaze(cell domain.Cell) []domain.CoordPair {
	neighbors := make([]domain.CoordPair, 0)
	cRow := cell.Row
	cCol := cell.Col

	if cell.Walls&domain.Up == 0 {
		neighbors = append(neighbors, domain.CoordPair{cRow - 1, cCol})
	}
	if cell.Walls&domain.Down == 0 {
		neighbors = append(neighbors, domain.CoordPair{cRow + 1, cCol})
	}
	if cell.Walls&domain.Left == 0 {
		neighbors = append(neighbors, domain.CoordPair{cRow, cCol - 1})
	}
	if cell.Walls&domain.Right == 0 {
		neighbors = append(neighbors, domain.CoordPair{cRow, cCol + 1})
	}
	return neighbors
}

func InitializeDistance(size int) []int {
	distance := make([]int, size)
	for i := range distance {
		distance[i] = infiniteWeight
	}
	distance[0] = 0
	return distance
}

func ChooseVertexWithShortestDistance(distance []int, visited []bool) int {
	minDistIdx := 0
	minValue := infiniteWeight
	for i := range distance {
		if distance[i] < minValue && !visited[i] {
			minValue = distance[i]
			minDistIdx = i
		}
	}
	return minDistIdx
}

func EncodeCoordinates(cord domain.CoordPair, width int) int {
	return cord[0]*width + cord[1]
}

func (d DijkstraSolver) SolveMaze(graph domain.Graph, startCoord, goalCoord domain.Coordinates) ([]int, error) {
	edges := TransformGraphIntoEdgeList(graph)
	visited := make([]bool, graph.Width*graph.Height)
	start := EncodeCoordinates(domain.CoordPair{startCoord.Row, startCoord.Col}, graph.Width)
	parent := make([]int, graph.Width*graph.Height)
	parent[start] = -1

	adj := make([]domain.CoordPair, 0, len(edges)*2)
	for _, edge := range edges {
		src := EncodeCoordinates(edge.NodeSource, graph.Width)
		dst := EncodeCoordinates(edge.NodeDest, graph.Width)
		adj = append(adj, domain.CoordPair{src, dst}, domain.CoordPair{dst, src})
	}

	distance := InitializeDistance(graph.Width * graph.Height)
	for !AllVerticesAreVisited(visited) {
		vertex := ChooseVertexWithShortestDistance(distance, visited)
		for _, e := range adj {
			source := e[0]
			v := e[1]
			if vertex != source || visited[v] {
				continue
			}
			sourceCoordRow := DecodeCoordinates(source, graph.Width)[0]
			sourceCoordCol := DecodeCoordinates(source, graph.Width)[1]
			sourceWeight := graph.Graph[sourceCoordRow][sourceCoordCol].Weight

			if distance[v] > distance[source]+sourceWeight {
				distance[v] = distance[source] + sourceWeight
				parent[v] = source
			}
		}
		visited[vertex] = true
	}
	return parent, nil
}

func AllVerticesAreVisited(visited []bool) bool {
	for i := range visited {
		if !visited[i] {
			return false
		}
	}
	return true
}
