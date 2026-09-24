package application

import "math"

type Fisheye struct {
	Weight float64
}

func NewFisheye(weight float64) Fisheye {
	return Fisheye{Weight: weight}
}

func (f Fisheye) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)

	xTransform := 2 * x / (r + 1)
	yTransform := 2 * y / (r + 1)
	return xTransform, yTransform
}

func (f Fisheye) GetWeight() float64 {
	return f.Weight
}
