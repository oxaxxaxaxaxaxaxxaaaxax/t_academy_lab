package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolar(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "polar (0,0)",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 0,
			expectedY: -1,
		},
		{
			name:      "polar (1,0)",
			x:         1,
			y:         0,
			weight:    0.5,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "polar (0,1)",
			x:         0,
			y:         1,
			weight:    2,
			expectedX: 0.5,
			expectedY: 0,
		},
		{
			name:      "polar (3,4)",
			x:         3,
			y:         4,
			weight:    1,
			expectedX: math.Atan2(4, 3) / math.Pi,
			expectedY: 5 - 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewPolar(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
