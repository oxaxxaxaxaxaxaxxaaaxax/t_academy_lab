package adapter

import (
	"testing"

	"hw4-fractal-flame/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewRender(t *testing.T) {
	tests := []struct {
		width, height  int
		expectedWidth  int
		expectedHeight int
	}{
		{
			width:          3,
			height:         2,
			expectedWidth:  3,
			expectedHeight: 2,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			cfg := domain.FractalConfig{Width: test.width, Height: test.height}

			r := NewRender(cfg)

			b := r.Image.Bounds()
			assert.Equal(t, test.expectedWidth, b.Dx())
			assert.Equal(t, test.expectedHeight, b.Dy())
		})
	}
}
