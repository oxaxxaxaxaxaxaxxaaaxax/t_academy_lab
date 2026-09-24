package application

import (
	"hw4-fractal-flame/internal/domain"
	"log/slog"
)

type Color struct {
	R, G, B uint8
}

func InitColorsArray(n int, rand domain.Random) []Color {
	slog.Debug("init array colors")

	colors := make([]Color, n)
	for i := range n {
		colors[i].R = uint8(rand.GetRandomIntInterval(64, 255))
		colors[i].G = uint8(rand.GetRandomIntInterval(64, 255))
		colors[i].B = uint8(rand.GetRandomIntInterval(64, 255))
	}
	return colors
}
