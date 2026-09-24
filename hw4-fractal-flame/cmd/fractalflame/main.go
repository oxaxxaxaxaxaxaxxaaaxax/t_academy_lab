package main

import (
	"hw4-fractal-flame/internal/adapter"
	"hw4-fractal-flame/internal/application"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	))
	cfg, err := adapter.ParseUserInput()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = application.StartChaosGame(cfg)
	if err != nil {
		slog.Error(err.Error())
	}
}
