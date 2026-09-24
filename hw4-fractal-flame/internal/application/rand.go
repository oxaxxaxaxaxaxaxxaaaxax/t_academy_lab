package application

import (
	"math/rand"
)

type ChaosRand struct {
	Rng rand.Rand
}

func NewChaosRand(seed int64) ChaosRand {
	return ChaosRand{Rng: *rand.New(rand.NewSource(seed))}
}

// return random float64 from [-module; module]
func (c ChaosRand) GetRandomFloat(module float64) float64 {
	f := float64(c.Rng.Intn(1001)) / 1000.00
	return f*(2*module) - module
}

// return random float64 from [0; module]
func (c ChaosRand) GetRandomFloatPositive(module float64) float64 {
	f := float64(c.Rng.Intn(1001)) / 1000.00
	return f * module
}

// // return random float64 from [0; module)
func (c ChaosRand) Float64n(module float64) float64 {
	return c.Rng.Float64() * module
}

// return random int from [0; module]
func (c ChaosRand) GetRandomIntPositive(module int) int {
	return c.Rng.Intn(module + 1)
}

// return random int from [to; from]
func (c ChaosRand) GetRandomIntInterval(from, to int) int {
	return c.Rng.Intn(to-from+1) + from
}
