package application

import "math"

type Cosine struct {
	Weight float64
}

func NewCosine(weight float64) Cosine {
	return Cosine{Weight: weight}
}

func (c Cosine) ApplyTransform(x, y float64) (float64, float64) {
	xTransform := math.Cos(math.Pi*x) * math.Cosh(y)
	yTransform := -math.Sin(math.Pi*x) * math.Sinh(y)
	return xTransform, yTransform
}

func (c Cosine) GetWeight() float64 {
	return c.Weight
}
