package main

import (
	"handman/internal/adapter"
	"handman/internal/application"
	"handman/internal/domain"
	"os"
)

func main() {
	ui := adapter.NewConsoleUI(os.Stdin)
	rnd := domain.NewWordRand()
	game := application.NewGame(ui, rnd)
	game.StartGame()
}
