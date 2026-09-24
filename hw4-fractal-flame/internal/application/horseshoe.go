package application

import "math"

type Horseshoe struct {
	Weight float64
}

func NewHorseshoe(weight float64) Horseshoe {
	return Horseshoe{Weight: weight}
}

func (h Horseshoe) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)

	xTransform := ((x - y) * (x + y)) / r
	yTransform := (2 * x * y) / r
	return xTransform, yTransform
}

func (h Horseshoe) GetWeight() float64 {
	return h.Weight
}
