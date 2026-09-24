package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const delta = 0.000000001

func TestSinusoidal(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "sinusoidal (zero values)",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "sinusoidal",
			x:         math.Pi / 2,
			y:         math.Pi,
			weight:    0.7,
			expectedX: 1,
			expectedY: 0,
		},
		{
			name:      "sinusoidal (negative values)",
			x:         -math.Pi / 2,
			y:         -math.Pi,
			weight:    -2,
			expectedX: -1,
			expectedY: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewSinusoidal(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
