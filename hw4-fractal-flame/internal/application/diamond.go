package application

import "math"

type Diamond struct {
	Weight float64
}

func NewDiamond(weight float64) Diamond {
	return Diamond{Weight: weight}
}

func (d Diamond) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(x, y)

	xTransform := math.Sin(t) * math.Cos(r)
	yTransform := math.Cos(t) * math.Sin(r)

	return xTransform, yTransform
}

func (d Diamond) GetWeight() float64 {
	return d.Weight
}
