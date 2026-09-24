package application

import "math"

type Bubble struct {
	Weight float64
}

func NewBubble(weight float64) Bubble {
	return Bubble{Weight: weight}
}

func (c Bubble) ApplyTransform(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	coef := 4 / ((r * r) + 4)
	xTransform := x * coef
	yTransform := y * coef
	return xTransform, yTransform
}

func (c Bubble) GetWeight() float64 {
	return c.Weight
}
