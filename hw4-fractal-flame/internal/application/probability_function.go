package application

import (
	"hw4-fractal-flame/internal/domain"
	"log/slog"
	"math"
	"sync"
)

type ProbabilityIteratedFunction struct {
	WeightSum  float64
	Intervals  []domain.Interval
	Variations []domain.Variation
}

func NewProbabilityIteratedFunction(variations []domain.Variation) ProbabilityIteratedFunction {
	return ProbabilityIteratedFunction{Variations: variations}
}

func (p ProbabilityIteratedFunction) InitWeightIntervals() ProbabilityIteratedFunction {
	oldWeightSum, newWeightSum := 0.0, 0.0

	for _, v := range p.Variations {
		newWeightSum += v.GetWeight()
		p.Intervals = append(p.Intervals, domain.Interval{From: oldWeightSum, To: newWeightSum})
		oldWeightSum = newWeightSum
	}

	p.WeightSum = newWeightSum
	return p
}

func (p ProbabilityIteratedFunction) RunIterations(cfg domain.FractalConfig) (domain.Pixels, error) {
	slog.Info("start RunIterations")

	canvas := domain.NewCanvas(cfg.Width, cfg.Height)
	affine, err := ParseAffineParams(cfg)
	if err != nil {
		return domain.Pixels{}, err
	}

	rand := NewChaosRand(cfg.Seed)
	colors := InitColorsArray(affine.N, rand)

	p = p.InitWeightIntervals()

	wg := sync.WaitGroup{}
	errChan := make(chan error, cfg.Threads)

	for thread := range cfg.Threads {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err = p.RunWorker(thread, cfg, canvas, affine, colors)
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

func (p ProbabilityIteratedFunction) TransformCoordinates(affineCoefficient AffineParams, x, y float64, rand ChaosRand, variations []domain.Variation) (float64, float64) {
	newX := affineCoefficient.A*x + affineCoefficient.B*y + affineCoefficient.C
	newY := affineCoefficient.D*x + affineCoefficient.E*y + affineCoefficient.F

	Vj := p.GetRandomVariations(rand, variations)
	return Vj.ApplyTransform(newX, newY)
}

func (p ProbabilityIteratedFunction) RunWorker(threadId int, cfg domain.FractalConfig,
	canvas domain.Canvas, affine CoefficientArray, colors []Color) error {
	iterCount := cfg.IterCount
	eqCount := affine.N
	pixels := canvas.Pixels
	rand := NewChaosRand(cfg.Seed + int64(threadId))
	variations := p.Variations

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

		for step := -20; step < iterCount; step++ {
			i := rand.GetRandomIntPositive(eqCount - 1)

			newX, newY = p.TransformCoordinates(affine.AffineCoefficients[i], newX, newY, rand, variations)

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

func (p ProbabilityIteratedFunction) GetRandomVariations(rand domain.Random, variations []domain.Variation) domain.Variation {
	weight := rand.Float64n(p.WeightSum)

	for idx, interval := range p.Intervals {
		if IsInInterval(interval, weight) {
			return variations[idx]
		}
	}
	return variations[0]
}

func IsInInterval(interval domain.Interval, weight float64) bool {
	return (interval.From <= weight) && (weight < interval.To)
}
