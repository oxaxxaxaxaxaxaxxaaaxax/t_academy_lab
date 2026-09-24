package application

import "math"

type Heart struct {
	Weight float64
}

func NewHeart(weight float64) Heart {
	return Heart{Weight: weight}
}

func (h Heart) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)

	xTransform := r * math.Sin(t*r)
	yTransform := -r * math.Cos(t*r)
	return xTransform, yTransform
}

func (h Heart) GetWeight() float64 {
	return h.Weight
}
