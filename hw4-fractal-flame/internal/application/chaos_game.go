package application

import (
	"hw4-fractal-flame/internal/adapter"
	"hw4-fractal-flame/internal/domain"
	"log/slog"
	"math"
)

const (
	samples = 10000
	xMax    = 1.777
	xMin    = -1.777
	yMax    = 1
	yMin    = -1
)

func StartChaosGame(cfg domain.FractalConfig) error {
	slog.Info("Starting chaos game")

	variations, err := InitVariations(cfg)
	if err != nil {
		return err
	}

	iter := NewIteratedFunction(cfg.IterAlgorithmName, variations)

	pixels, err := iter.RunIterations(cfg)
	if err != nil {
		return err
	}

	pixels = GammaCorrection(pixels, cfg.GammaCoefficient)

	render := adapter.NewRender(cfg)
	err = render.RenderImage(pixels, cfg.OutputPath)
	if err != nil {
		return err
	}

	return nil
}

func GammaCorrection(pixels domain.Pixels, gamma float64) domain.Pixels {
	maxV := 0.0
	gammaPower := 1.0 / gamma

	for rowIdx, row := range pixels {
		for colIdx := range row {
			p := &pixels[rowIdx][colIdx]

			if p.HitCounter == 0 {
				continue
			}

			p.NormalizedValue = math.Log10(float64(p.HitCounter))
			if p.NormalizedValue > maxV {
				maxV = p.NormalizedValue
			}
		}
	}

	for rowIdx, row := range pixels {
		for colIdx, _ := range row {
			pixels[rowIdx][colIdx].NormalizedValue /= maxV
			gammaCoef := math.Pow(pixels[rowIdx][colIdx].NormalizedValue, gammaPower)

			pixels[rowIdx][colIdx].R = CutToU8(float64(pixels[rowIdx][colIdx].R) * gammaCoef)
			pixels[rowIdx][colIdx].G = CutToU8(float64(pixels[rowIdx][colIdx].G) * gammaCoef)
			pixels[rowIdx][colIdx].B = CutToU8(float64(pixels[rowIdx][colIdx].B) * gammaCoef)
		}
	}

	return pixels
}

func CutToU8(color float64) uint8 {
	if color > 255 {
		return 255
	}
	if color < 0 {
		return 0
	}
	return uint8(color)
}

func InitVariations(cfg domain.FractalConfig) ([]domain.Variation, error) {
	var variations []domain.Variation

	names, err := adapter.ParseVariations(cfg)
	if err != nil {
		return nil, err
	}

	for name, weight := range names {
		variations = append(variations, NewVariation(name, weight))
	}

	return variations, nil
}
