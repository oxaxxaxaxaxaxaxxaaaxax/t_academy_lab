package application

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"

type Generator interface {
	GenerateMaze(height, width int) (domain.Graph, error)
}

type Solver interface {
	SolveMaze(graph domain.Graph, start, goal domain.Coordinates) ([]int, error)
}

type Drawer interface {
	ShowMaze(graph domain.Graph)
	DrawSolution(graph domain.Graph, path []int, start, goal domain.Coordinates)
}
