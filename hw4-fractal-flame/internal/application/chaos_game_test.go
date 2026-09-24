package application

import (
	"hw4-fractal-flame/internal/adapter"
	"testing"

	"hw4-fractal-flame/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestIsInInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval domain.Interval
		weight   float64
		expected bool
	}{
		{
			name:     "inside interval",
			interval: domain.Interval{From: 0.0, To: 1.0},
			weight:   0.5,
			expected: true,
		},
		{
			name:     "from including",
			interval: domain.Interval{From: 0.3, To: 0.7},
			weight:   0.3,
			expected: true,
		},
		{
			name:     "to excluding",
			interval: domain.Interval{From: 0.3, To: 0.7},
			weight:   0.7,
			expected: false,
		},
		{
			name:     "outside interval",
			interval: domain.Interval{From: 0.3, To: 0.7},
			weight:   0.2,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := IsInInterval(test.interval, test.weight)
			assert.Equal(t, test.expected, in)
		})
	}
}

func TestCutToU8(t *testing.T) {
	tests := []struct {
		name     string
		color    float64
		expected uint8
	}{
		{name: "color < 0", color: -1, expected: 0},
		{name: "color == 0", color: 0, expected: 0},
		{name: "0 < color < 255", color: 12.9, expected: 12},
		{name: "color == 255", color: 255, expected: 255},
		{name: "color > 255", color: 9999, expected: 255},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			color := CutToU8(test.color)
			assert.Equal(t, test.expected, color)
		})
	}
}

func TestGammaCorrection(t *testing.T) {
	tests := []struct {
		name           string
		pixels         domain.Pixels
		gamma          float64
		expectedPixels domain.Pixels
	}{
		{
			name: "gamma=1, hitcounter != 0",
			pixels: domain.Pixels{
				{
					{R: 100, G: 100, B: 100, HitCounter: 10},
					{R: 100, G: 100, B: 100, HitCounter: 100},
				},
			},
			gamma: 1.0,
			expectedPixels: domain.Pixels{
				{
					{R: 50, G: 50, B: 50, HitCounter: 10, NormalizedValue: 0.5},
					{R: 100, G: 100, B: 100, HitCounter: 100, NormalizedValue: 1.0},
				},
			},
		},
		{
			name: "gamma=1, hitcounter == 0",
			pixels: domain.Pixels{
				{
					{R: 200, G: 10, B: 10, HitCounter: 0},
					{R: 100, G: 100, B: 100, HitCounter: 100},
				},
			},
			gamma: 1.0,
			expectedPixels: domain.Pixels{
				{
					{R: 0, G: 0, B: 0, HitCounter: 0, NormalizedValue: 0},
					{R: 100, G: 100, B: 100, HitCounter: 100, NormalizedValue: 1.0},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testColors := GammaCorrection(test.pixels, test.gamma)

			assert.Equal(t, test.expectedPixels[0][0].R, testColors[0][0].R)
			assert.Equal(t, test.expectedPixels[0][0].G, testColors[0][0].G)
			assert.Equal(t, test.expectedPixels[0][0].B, testColors[0][0].B)

			assert.Equal(t, test.expectedPixels[0][1].R, testColors[0][1].R)
			assert.Equal(t, test.expectedPixels[0][1].G, testColors[0][1].G)
			assert.Equal(t, test.expectedPixels[0][1].B, testColors[0][1].B)

			assert.Equal(t, test.expectedPixels[0][0].NormalizedValue, testColors[0][0].NormalizedValue)
			assert.Equal(t, test.expectedPixels[0][1].NormalizedValue, testColors[0][1].NormalizedValue)
		})
	}
}

func TestInitVariations(t *testing.T) {
	tests := []struct {
		name            string
		functions       string
		expectedLen     int
		expectedErr     error
		expectedWeights []float64
	}{
		{
			name:            "single function",
			functions:       "heart:0.7",
			expectedLen:     1,
			expectedWeights: []float64{0.7},
			expectedErr:     nil,
		},
		{
			name:            "multiple functions",
			functions:       "heart:0.7,linear:0.3",
			expectedLen:     2,
			expectedWeights: []float64{0.7, 0.3},
			expectedErr:     nil,
		},
		{
			name:        "invalid weight",
			functions:   "heart:AAA",
			expectedErr: adapter.ErrInvalidVariations,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := domain.FractalConfig{Functions: test.functions}

			variations, err := InitVariations(cfg)

			assert.ErrorIs(t, test.expectedErr, err)
			assert.Len(t, variations, test.expectedLen)

			var variationWeights []float64
			for _, v := range variations {
				variationWeights = append(variationWeights, v.GetWeight())
			}
			assert.ElementsMatch(t, test.expectedWeights, variationWeights)
		})
	}
}
