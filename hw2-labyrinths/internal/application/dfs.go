package application

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"

type DFSGenerator struct{}

func NewDFSGenerator() DFSGenerator {
	return DFSGenerator{}
}

func (d DFSGenerator) InitializeStart(graph [][]domain.Cell) ([][]domain.Cell, domain.Stack) {
	graph[0][0].IsVisited = true
	st := domain.Stack{St: make([]domain.Cell, 0)}
	st = st.Push(graph[0][0])
	return graph, st
}

func (d DFSGenerator) GenerateMaze(height, width int) (domain.Graph, error) {
	initialGraph := CreateGraph(height, width)
	graph := initialGraph.Graph

	graph, st := d.InitializeStart(graph)

	for !st.Empty() {
		stCopy, cell, err := st.Pop()

		if err != nil {
			return domain.Graph{}, err
		}

		neighbors := getNeighbors(width, height, cell.Row, cell.Col)
		if len(neighbors) != 0 {
			for _, neighbor := range neighbors {
				nr := neighbor[0]
				nc := neighbor[1]
				if !graph[nr][nc].IsVisited {
					cell, graph[nr][nc] = DoPassage(cell, graph[nr][nc])
					stCopy = stCopy.Push(cell)
					graph[cell.Row][cell.Col] = cell
					graph[nr][nc].IsVisited = true
					stCopy = stCopy.Push(graph[nr][nc])
					st = stCopy
					break
				}
			}
		}
		st = stCopy

	}
	initialGraph.Graph = graph
	return initialGraph, nil
}
