package application

import "hw4-fractal-flame/internal/domain"

func NewIteratedFunction(name string, variations []domain.Variation) domain.IteratedFunction {
	switch name {
	case "sum_iter":
		return NewSumIteratedFunction(variations)
	case "prob_iter":
		return NewProbabilityIteratedFunction(variations)
	default:
		return NewSumIteratedFunction(variations)
	}
}

func NewVariation(name string, weight float64) domain.Variation {
	switch name {
	case "linear":
		return NewLinear(weight)
	case "horseshoe":
		return NewHorseshoe(weight)
	case "sinusoidal":
		return NewSinusoidal(weight)
	case "heart":
		return NewHeart(weight)
	case "fisheye":
		return NewFisheye(weight)
	case "diamond":
		return NewDiamond(weight)
	case "cosine":
		return NewCosine(weight)
	case "swirl":
		return NewSwirl(weight)
	case "polar":
		return NewPolar(weight)
	case "bubble":
		return NewBubble(weight)
	case "spherical":
		return NewSpherical(weight)
	}
	return nil
}
