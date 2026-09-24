package application

import (
	"log/slog"
	"math"
	"sync"

	"hw4-fractal-flame/internal/domain"
)

const defaultIterCount = 0

type SumIteratedFunction struct {
	Variations []domain.Variation
}

func NewSumIteratedFunction(variations []domain.Variation) SumIteratedFunction {
	return SumIteratedFunction{Variations: variations}
}

func (s SumIteratedFunction) InitVariations(variations []domain.Variation) SumIteratedFunction {
	return SumIteratedFunction{Variations: variations}
}

func (s SumIteratedFunction) RunIterations(cfg domain.FractalConfig) (domain.Pixels, error) {
	slog.Info("start RunIterations")

	canvas := domain.NewCanvas(cfg.Width, cfg.Height)
	affine, err := ParseAffineParams(cfg)
	if err != nil {
		return domain.Pixels{}, err
	}

	rand := NewChaosRand(cfg.Seed)
	colors := InitColorsArray(affine.N, rand)

	wg := sync.WaitGroup{}
	errChan := make(chan error, cfg.Threads)

	for thread := range cfg.Threads {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err = s.RunWorker(thread, cfg, canvas, affine, colors)
			if err != nil {
				slog.Error(err.Error())
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err = range errChan {
		return domain.Pixels{}, err
	}
	return canvas.Pixels, nil
}

func (s SumIteratedFunction) TransformCoordinates(affineCoefficient AffineParams, x, y float64, variations []domain.Variation) (float64, float64) {
	newX := affineCoefficient.A*x + affineCoefficient.B*y + affineCoefficient.C
	newY := affineCoefficient.D*x + affineCoefficient.E*y + affineCoefficient.F

	Fx, Fy := 0.0, 0.0
	for idx, v := range variations {
		Vj := variations[idx]

		xSumCoef, ySumCoef := Vj.ApplyTransform(newX, newY)
		Fx += xSumCoef * v.GetWeight()
		Fy += ySumCoef * v.GetWeight()
	}
	return Fx, Fy
}

func (s SumIteratedFunction) RunWorker(threadId int, cfg domain.FractalConfig,
	canvas domain.Canvas, affine CoefficientArray, colors []Color) error {
	iterCount := cfg.IterCount
	eqCount := affine.N
	pixels := canvas.Pixels
	rand := NewChaosRand(cfg.Seed + int64(threadId))
	variations := s.Variations

	symmetryLevel := cfg.SymmetryLevel
	workerSamples := samples / cfg.Threads

	transformX := float64(cfg.Width-1) / (xMax - xMin)
	transformY := float64(cfg.Height-1) / (yMax - yMin)

	for c := range workerSamples {
		newX := rand.GetRandomFloat(xMax)
		newY := rand.GetRandomFloat(yMax)

		if c%1000 == 0 {
			slog.Debug("samples count", "count:", c)
		}

		for step := defaultIterCount; step < iterCount; step++ {
			i := rand.GetRandomIntPositive(eqCount - 1)

			newX, newY = s.TransformCoordinates(affine.AffineCoefficients[i], newX, newY, variations)

			if step <= 0 {
				continue
			}

			theta := 0.0
			for range symmetryLevel {
				theta += (2 * math.Pi) / float64(symmetryLevel)
				xRot := newX*math.Cos(theta) - newY*math.Sin(theta)
				yRot := newX*math.Sin(theta) + newY*math.Cos(theta)

				if xRot <= xMax && xRot >= xMin && yRot <= yMax && yRot >= yMin {
					coordX := int((xRot - xMin) * transformX)
					coordY := int((yRot - yMin) * transformY)

					canvas.Mutexes[coordY].Lock()

					point := pixels[coordY][coordX]
					if point.HitCounter == 0 {
						point.R = colors[i].R
						point.G = colors[i].G
						point.B = colors[i].B
					} else {
						point.R = (point.R + colors[i].R) / 2
						point.G = (point.G + colors[i].G) / 2
						point.B = (point.B + colors[i].B) / 2
					}

					point.HitCounter++
					pixels[coordY][coordX] = point

					canvas.Mutexes[coordY].Unlock()
				}
			}
		}
	}
	return nil
}
