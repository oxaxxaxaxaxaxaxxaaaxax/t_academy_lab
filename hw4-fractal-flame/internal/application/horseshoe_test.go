package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHorseshoe(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "horseshoe (1,0)",
			x:         1,
			y:         0,
			weight:    1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			name:      "horseshoe (0,1)",
			x:         0,
			y:         1,
			weight:    0.5,
			expectedX: -1,
			expectedY: 0,
		},
		{
			name:      "horseshoe (3,4)",
			x:         3,
			y:         4,
			weight:    -2,
			expectedX: -7.0 / 5.0,
			expectedY: 24.0 / 5.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewHorseshoe(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
