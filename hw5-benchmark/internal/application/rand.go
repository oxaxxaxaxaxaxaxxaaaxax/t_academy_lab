package application

import (
	"math/rand"
	"strings"
)

const (
	intMax       = 1000000
	floatMax     = 1000000.0
	maxStringLen = 50
)

type InspectorRand struct {
	rng rand.Rand
}

func NewInspectorRand() InspectorRand {
	return InspectorRand{rng: *rand.New(rand.NewSource(42))}
}

func (i InspectorRand) RandBoll() bool {
	return i.rng.Intn(2) == 0
}

func (i InspectorRand) RandInt64() int64 {
	return i.rng.Int63n(2*intMax) - intMax
}

func (i InspectorRand) RandIntPositive() int {
	return i.rng.Intn(intMax)
}

func (i InspectorRand) RandUint64() uint64 {
	return i.rng.Uint64()
}

func (i InspectorRand) RandFloat64() float64 {
	f := float64(i.rng.Intn(1001)) / 1000.00
	return f*(2*floatMax) - floatMax
}

func (i InspectorRand) RandString() string {
	var b strings.Builder
	stringLen := i.rng.Intn(maxStringLen)

	for range stringLen {
		b.WriteString(string(i.rng.Int31()))
	}

	return b.String()
}
