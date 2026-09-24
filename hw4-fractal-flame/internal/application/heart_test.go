package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeart(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "heart (0,0)",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "heart (1,0)",
			x:         1,
			y:         0,
			weight:    0.7,
			expectedX: 0,
			expectedY: -1,
		},
		{
			name:      "heart (0,1)",
			x:         0,
			y:         1,
			weight:    -2,
			expectedX: math.Sin(math.Pi/2) * 1,
			expectedY: -math.Cos(math.Pi/2) * 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewHeart(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
