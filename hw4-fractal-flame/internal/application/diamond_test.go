package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiamond(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "diamond (0,0)",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "diamond (0,1)",
			x:         0,
			y:         1,
			weight:    0.7,
			expectedX: 0,
			expectedY: math.Sin(1),
		},
		{
			name:      "diamond (3,4)",
			x:         3,
			y:         4,
			weight:    -2,
			expectedX: math.Sin(math.Atan2(3, 4)) * math.Cos(5),
			expectedY: math.Cos(math.Atan2(3, 4)) * math.Sin(5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewDiamond(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
