package application

import (
	"fmt"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type SolverService struct {
	solver Solver
	drawer Drawer
}

func NewSolverService(solver Solver, drawer Drawer) SolverService {
	return SolverService{solver: solver, drawer: drawer}
}

func (s SolverService) SolveMaze(graph domain.Graph, start, goal domain.Coordinates) {
	path, err := s.solver.SolveMaze(graph, start, goal)
	if err != nil {
		adapter.ShowError(err)
		return
	}
	fmt.Println(path)
	s.drawer.DrawSolution(graph, path, start, goal)
}
