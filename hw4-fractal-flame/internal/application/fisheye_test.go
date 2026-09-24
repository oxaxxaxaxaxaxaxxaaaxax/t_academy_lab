package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFisheye(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "fisheye (0,0)",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "fisheye (1,0)",
			x:         1,
			y:         0,
			weight:    0.7,
			expectedX: 1,
			expectedY: 0,
		},
		{
			name:      "fisheye (3,4)",
			x:         3,
			y:         4,
			weight:    -2,
			expectedX: 6.0 / (5.0 + 1.0),
			expectedY: 8.0 / (5.0 + 1.0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewFisheye(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
