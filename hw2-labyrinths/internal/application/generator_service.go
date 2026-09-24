package application

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"

type GeneratorService struct {
	generator Generator
	drawer    Drawer
}

func NewGeneratorService(generator Generator, drawer Drawer) GeneratorService {
	return GeneratorService{generator: generator, drawer: drawer}
}

func (g GeneratorService) CreateMaze(height, width int) {
	graph, err := g.generator.GenerateMaze(height, width)
	if err != nil {
		adapter.ShowError(err)
		return
	}
	g.drawer.ShowMaze(graph)
	return
}
