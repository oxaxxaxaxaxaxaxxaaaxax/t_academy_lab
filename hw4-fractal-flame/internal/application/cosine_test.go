package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosine(t *testing.T) {
	tests := []struct {
		x, y      float64
		weight    float64
		name      string
		expectedX float64
		expectedY float64
	}{
		{
			name:      "cosine (0,0)",
			x:         0,
			y:         0,
			weight:    1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			name:      "cosine (0.5,0)",
			x:         0.5,
			y:         0,
			weight:    0.7,
			expectedX: math.Cos(math.Pi*0.5) * math.Cosh(0),
			expectedY: -math.Sin(math.Pi*0.5) * math.Sinh(0),
		},
		{
			name:      "cosine (1,1)",
			x:         1,
			y:         1,
			weight:    -2,
			expectedX: math.Cos(math.Pi*1) * math.Cosh(1),
			expectedY: -math.Sin(math.Pi*1) * math.Sinh(1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewCosine(test.weight)

			assert.Equal(t, test.weight, v.GetWeight())

			newX, newY := v.ApplyTransform(test.x, test.y)

			assert.InDelta(t, test.expectedX, newX, delta)
			assert.InDelta(t, test.expectedY, newY, delta)
		})
	}
}
