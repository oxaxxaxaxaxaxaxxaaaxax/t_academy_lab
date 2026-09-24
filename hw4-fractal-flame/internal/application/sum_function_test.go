package application

import (
	"fmt"
	"hw4-fractal-flame/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkSumIteratedFunction_RunIterations(b *testing.B) {
	threadCases := []int{1, 2, 4, 8}

	for _, threads := range threadCases {
		cfg := domain.NewFractalConfig(
			920, 580, 50000, threads, 6,
			3, "bench.png",
			"0.5,1.0,1.0,1.0,1.0,1.0",
			"heart:0.7",
			"", "sum_iter", 0.7,
		)
		variations, err := InitVariations(cfg)
		if err != nil {
			b.Fatal(err.Error())
		}
		b.Run(fmt.Sprintf("threads=%d", threads), func(b *testing.B) {
			s := NewSumIteratedFunction(variations)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = s.RunIterations(cfg)
			}
		})

	}
}

func TestSumIteratedFunction_RunIterations(t *testing.T) {
	cfg := domain.FractalConfig{
		Width:         10,
		Height:        10,
		IterCount:     20,
		Threads:       1,
		Seed:          42,
		SymmetryLevel: 1,
		AffineParams:  "0.5,0,0,0,0.5,0",
		Functions:     "linear:1.0",
	}

	variations := []domain.Variation{
		NewLinear(1.0),
	}

	s := NewSumIteratedFunction(variations)

	pixels, err := s.RunIterations(cfg)
	require.NoError(t, err)

	require.Len(t, pixels, cfg.Height)
	require.Len(t, pixels[0], cfg.Width)

	hitCount := false
	for _, row := range pixels {
		for _, p := range row {
			if p.HitCounter > 0 {
				hitCount = true
				break
			}
		}
		if hitCount {
			break
		}
	}

	assert.True(t, hitCount)
}
