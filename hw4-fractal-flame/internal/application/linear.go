package application

type Linear struct {
	Weight float64
}

func NewLinear(weight float64) Linear {
	return Linear{Weight: weight}
}

func (l Linear) ApplyTransform(x, y float64) (float64, float64) {
	return x, y
}

func (l Linear) GetWeight() float64 {
	return l.Weight
}
