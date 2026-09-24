package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinear(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "linear coordinates",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "linear coordinates (positive weight)",
			x:         2.5,
			y:         -3.75,
			weight:    0.7,
			expectedX: 2.5,
			expectedY: -3.75,
		},
		{
			name:      "linear coordinates (negative weight)",
			x:         -10,
			y:         4,
			weight:    -2,
			expectedX: -10,
			expectedY: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewLinear(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.Equal(t, test.expectedX, newX)
			assert.Equal(t, test.expectedY, newY)
		})
	}
}
