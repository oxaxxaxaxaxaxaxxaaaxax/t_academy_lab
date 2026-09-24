package application

import "math"

type Sinusoidal struct {
	Weight float64
}

func NewSinusoidal(weight float64) Sinusoidal {
	return Sinusoidal{Weight: weight}
}

func (s Sinusoidal) ApplyTransform(x, y float64) (float64, float64) {
	xTransform := math.Sin(x)
	yTransform := math.Sin(y)
	return xTransform, yTransform
}

func (s Sinusoidal) GetWeight() float64 {
	return s.Weight
}
