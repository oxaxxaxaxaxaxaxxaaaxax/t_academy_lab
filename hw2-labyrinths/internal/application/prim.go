package application

import (
	"math/rand"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type PrimGenerator struct{}

func NewPrimGenerator() PrimGenerator {
	return PrimGenerator{}
}

func CreateGraph(height, width int) domain.Graph {
	graph := make([][]domain.Cell, height)
	for rowIdx := range height {
		graph[rowIdx] = make([]domain.Cell, width)
		for colIdx := range width {
			graph[rowIdx][colIdx] = CreateCellWithSpecialWeight(rowIdx, colIdx)
		}
	}
	return domain.Graph{Height: height, Width: width, Graph: graph}
}

// Create special cells by choosing x,y coordinates
func CreateCellWithSpecialWeight(rowIdx, colIdx int) domain.Cell {
	switch {
	case rowIdx == 0 && colIdx == 0:
		cellWeight := PassWeight
		return domain.NewCell(rowIdx, colIdx,
			domain.Up|domain.Down|domain.Left|domain.Right, cellWeight)
	case (colIdx%6 == 0 && rowIdx%4 == 0) || (colIdx%9 == 0 && rowIdx%7 == 0):
		cellWeight := CoinWeight
		return domain.NewCell(rowIdx, colIdx,
			domain.Up|domain.Down|domain.Left|domain.Right, cellWeight)
	case (colIdx%6 == 3 && rowIdx%5 == 0) || (colIdx%7 == 0 && rowIdx%7 == 0):
		cellWeight := SandWeight
		return domain.NewCell(rowIdx, colIdx,
			domain.Up|domain.Down|domain.Left|domain.Right, cellWeight)
	default:
		cellWeight := PassWeight
		return domain.NewCell(rowIdx, colIdx,
			domain.Up|domain.Down|domain.Left|domain.Right, cellWeight)
	}
}

func InGrid(width, height, row, col int) bool {
	return row < height && row >= 0 && col < width && col >= 0
}

func getNeighbors(width, height, row, col int) []domain.CoordPair {
	directions := []domain.CoordPair{{-1, 0}, {+1, 0}, {0, -1}, {0, +1}}

	neighbors := make([]domain.CoordPair, 0, 4)

	for _, val := range directions {
		rowNeighbor := row + val[0]
		colNeighbor := col + val[1]
		if InGrid(width, height, rowNeighbor, colNeighbor) {
			neighbors = append(neighbors, domain.CoordPair{rowNeighbor, colNeighbor})
		}
	}
	return neighbors
}

func (PrimGenerator) InitializeStart(graph [][]domain.Cell, height, width int) ([][]domain.Cell, []domain.Wall) {
	frontier := make([]domain.Wall, 0, 4)
	graph[0][0].IsVisited = true
	neighbors := getNeighbors(width, height, 0, 0)

	for _, neighbor := range neighbors {
		nr := neighbor[0]
		nc := neighbor[1]
		if !graph[nr][nc].IsVisited {
			frontier = append(frontier, domain.NewWall(graph[0][0], graph[nr][nc]))
		}
	}
	return graph, frontier
}

func DoPassage(From, To domain.Cell) (domain.Cell, domain.Cell) {
	rFrom := From.Row
	rTo := To.Row
	cFrom := From.Col
	cTo := To.Col
	rDelta := rTo - rFrom
	cDelta := cTo - cFrom
	switch {
	case rDelta == 0 && cDelta == -1:
		From.Walls = From.Walls & (^domain.Left)
		To.Walls = To.Walls & (^domain.Right)
	case rDelta == 0 && cDelta == 1:
		From.Walls = From.Walls & (^domain.Right)
		To.Walls = To.Walls & (^domain.Left)
	case rDelta == -1 && cDelta == 0:
		From.Walls = From.Walls & (^domain.Up)
		To.Walls = To.Walls & (^domain.Down)
	case rDelta == 1 && cDelta == 0:
		From.Walls = From.Walls & (^domain.Down)
		To.Walls = To.Walls & (^domain.Up)
	}

	return From, To
}

func (p PrimGenerator) GenerateMaze(height, width int) (domain.Graph, error) {
	initialGraph := CreateGraph(height, width)
	graph := initialGraph.Graph
	graph, frontier := p.InitializeStart(graph, height, width)
	if len(frontier) == 0 {
		initialGraph.Graph = graph
		return initialGraph, nil
	}

	for len(frontier) > 0 {
		idx := rand.Intn(len(frontier))
		f := frontier[idx]

		frontier[idx] = frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

		rowTo := f.CellToRow
		colTo := f.CellToCol
		rowFrom := f.CellFromRow
		colFrom := f.CellFromCol

		if graph[rowTo][colTo].IsVisited {
			continue
		}

		graph[rowFrom][colFrom], graph[rowTo][colTo] = DoPassage(graph[f.CellFromRow][f.CellFromCol],
			graph[f.CellToRow][f.CellToCol])

		graph[rowTo][colTo].IsVisited = true
		neighbors := getNeighbors(width, height, rowTo, colTo)

		for _, neighbor := range neighbors {
			nr := neighbor[0]
			nc := neighbor[1]
			if !graph[nr][nc].IsVisited {
				frontier = append(frontier, domain.NewWall(graph[rowTo][colTo], graph[nr][nc]))
			}
		}
	}
	initialGraph.Graph = graph
	return initialGraph, nil
}
