package application

import (
	"errors"
	"testing"

	"hw4-fractal-flame/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestParseAffineParams(t *testing.T) {
	tests := []struct {
		name           string
		affine         string
		expectedErr    error
		expectedN      int
		expectedCoeffs []AffineParams
	}{
		{
			name:      "single affine set",
			affine:    "0.5,1.0,1.0,1.0,1.0,1.0",
			expectedN: 1,
			expectedCoeffs: []AffineParams{
				{A: 0.5, B: 1.0, C: 1.0, D: 1.0, E: 1.0, F: 1.0},
			},
			expectedErr: nil,
		},
		{
			name:      "multiple affine sets",
			affine:    "0.48,0,0,0,0.48,0/0.26,0.16,0.28,-0.16,0.26,0",
			expectedN: 2,
			expectedCoeffs: []AffineParams{
				{A: 0.48, B: 0, C: 0, D: 0, E: 0.48, F: 0},
				{A: 0.26, B: 0.16, C: 0.28, D: -0.16, E: 0.26, F: 0},
			},
			expectedErr: nil,
		},
		{
			name:        "error not enough coefficient",
			affine:      "0.5,1.0,1.0,1.0,1.0",
			expectedErr: ErrInvalidConfig,
		},
		{
			name:        "error coefficient is not a number",
			affine:      "0.5,1.0,not_a_number,1.0,1.0,1.0",
			expectedErr: ErrInvalidConfig,
		},
		{
			name:        "error empty after slash",
			affine:      "0.5,1.0,1.0,1.0,1.0,1.0/",
			expectedErr: ErrInvalidConfig,
		},
		{
			name:        "error params string is empty",
			affine:      "",
			expectedErr: ErrInvalidConfig,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := domain.FractalConfig{
				AffineParams: test.affine,
			}

			result, err := ParseAffineParams(cfg)

			if test.expectedErr != nil {
				assert.True(t, errors.Is(err, test.expectedErr))
				assert.True(t, errors.Is(err, ErrWrongAffineParameters))
				return
			}

			assert.Equal(t, test.expectedN, result.N)
			assert.Equal(t, test.expectedCoeffs, result.AffineCoefficients)
		})
	}
}
