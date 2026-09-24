package application

import "math"

type Swirl struct {
	Weight float64
}

func NewSwirl(weight float64) Swirl {
	return Swirl{Weight: weight}
}

func (s Swirl) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)

	xTransform := x*math.Sin(r*r) - y*math.Cos(r*r)
	yTransform := x*math.Cos(r*r) - y*math.Sin(r*r)
	return xTransform, yTransform
}

func (s Swirl) GetWeight() float64 {
	return s.Weight
}
