package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const chaosSeed int64 = 42

func TestChaosRand_GetRandomFloat(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		module float64
	}{
		{name: "module=0", seed: chaosSeed, module: 0},
		{name: "module=1", seed: chaosSeed, module: 1},
		{name: "module=5", seed: chaosSeed, module: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewChaosRand(test.seed)

			for i := 0; i < 50; i++ {
				v := r.GetRandomFloat(test.module)
				assert.GreaterOrEqual(t, v, -test.module)
				assert.LessOrEqual(t, v, test.module)
			}
		})
	}
}

func TestChaosRand_GetRandomFloatPositive(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		module float64
	}{
		{name: "module=0", seed: chaosSeed, module: 0},
		{name: "module=1", seed: chaosSeed, module: 1},
		{name: "module=5", seed: chaosSeed, module: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewChaosRand(test.seed)

			for i := 0; i < 50; i++ {
				v := r.GetRandomFloatPositive(test.module)
				assert.GreaterOrEqual(t, v, 0.0)
				assert.LessOrEqual(t, v, test.module)
			}
		})
	}
}

func TestChaosRand_Float64n(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		module float64
	}{
		{name: "module=1 ", seed: chaosSeed, module: 1},
		{name: "module=100", seed: chaosSeed, module: 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewChaosRand(test.seed)

			for i := 0; i < 50; i++ {
				v := r.Float64n(test.module)
				assert.GreaterOrEqual(t, v, 0.0)
				assert.Less(t, v, test.module)
			}
		})
	}
}

func TestChaosRand_GetRandomIntPositive(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		module int
	}{
		{name: "module=0", seed: chaosSeed, module: 0},
		{name: "module=1", seed: chaosSeed, module: 1},
		{name: "module=10", seed: chaosSeed, module: 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewChaosRand(test.seed)

			for i := 0; i < 50; i++ {
				v := r.GetRandomIntPositive(test.module)
				assert.GreaterOrEqual(t, v, 0)
				assert.LessOrEqual(t, v, test.module)
			}
		})
	}
}

func TestChaosRand_GetRandomIntInterval(t *testing.T) {
	tests := []struct {
		name string
		seed int64
		from int
		to   int
	}{
		{name: "from=to", seed: chaosSeed, from: 5, to: 5},
		{name: "positive", seed: chaosSeed, from: 2, to: 8},
		{name: "module", seed: chaosSeed, from: -3, to: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewChaosRand(test.seed)

			for i := 0; i < 50; i++ {
				v := r.GetRandomIntInterval(test.from, test.to)
				assert.GreaterOrEqual(t, v, test.from)
				assert.LessOrEqual(t, v, test.to)
			}
		})
	}
}
