package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/application"
)

func main() {
	config, err := application.ParseCmdLine()
	if err != nil {
		adapter.ShowError(err)
		return
	}

	drawer, err := application.NewDrawer(config.Output, config.EnableUnicode)
	if err != nil {
		adapter.ShowError(err)
		return
	}

	switch {
	case config.Generate:
		generator := application.NewGenerator(config.Generator)
		service := application.NewGeneratorService(generator, drawer)
		service.CreateMaze(config.Height, config.Width)
	case config.Solve:
		solver := application.NewSolver(config.Solver)
		service := application.NewSolverService(solver, drawer)
		graph, mazeConfig, err := application.ParseMaze(config.File)
		if err != nil {
			adapter.ShowError(err)
			return
		}
		service.SolveMaze(graph, mazeConfig.Start, mazeConfig.Goal)
	default:
		//help message if user don't choose mode
		flag.Usage()
		_, err := fmt.Fprintf(os.Stdout, "current generate mode:%v ", config.Generate)
		if err != nil {
			adapter.ShowError(err)
			return
		}
		_, err = fmt.Fprintf(os.Stdout, "current solve mode:%v ", config.Solve)
		if err != nil {
			adapter.ShowError(err)
			return
		}
	}

}
