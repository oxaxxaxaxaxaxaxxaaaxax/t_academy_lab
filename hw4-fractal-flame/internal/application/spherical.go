package application

import "math"

type Spherical struct {
	Weight float64
}

func NewSpherical(weight float64) Spherical {
	return Spherical{Weight: weight}
}

func (s Spherical) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)

	xTransform := x / (r * r)
	yTransform := y / (r * r)
	return xTransform, yTransform
}

func (s Spherical) GetWeight() float64 {
	return s.Weight
}
