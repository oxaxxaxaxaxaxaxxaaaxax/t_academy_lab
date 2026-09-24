package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpherical(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "spherical (1,0)",
			x:         1,
			y:         0,
			weight:    1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			name:      "spherical (2,0)",
			x:         2,
			y:         0,
			weight:    0.5,
			expectedX: 0.5,
			expectedY: 0,
		},
		{
			name:      "spherical (3,4)",
			x:         3,
			y:         4,
			weight:    2,
			expectedX: 3.0 / 25.0,
			expectedY: 4.0 / 25.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewSpherical(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
