package domain

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks.go -package=mocks

type IteratedFunction interface {
	RunIterations(cfg FractalConfig) (Pixels, error)
}

type Random interface {
	GetRandomFloat(module float64) float64
	GetRandomFloatPositive(module float64) float64
	Float64n(module float64) float64
	GetRandomIntPositive(module int) int
	GetRandomIntInterval(from, to int) int
}

type Variation interface {
	ApplyTransform(x, y float64) (float64, float64)
	GetWeight() float64
}
