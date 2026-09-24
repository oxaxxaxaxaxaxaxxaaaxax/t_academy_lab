package application

import "math"

type Polar struct {
	Weight float64
}

func NewPolar(weight float64) Polar {
	return Polar{Weight: weight}
}

func (h Polar) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	return t / math.Pi, r - 1
}

func (h Polar) GetWeight() float64 {
	return h.Weight
}
