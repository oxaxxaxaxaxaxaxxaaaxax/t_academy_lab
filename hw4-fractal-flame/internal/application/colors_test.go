package application

import (
	"testing"

	"hw4-fractal-flame/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInitColorsArray(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		randomSequence []int
		expected       []Color
	}{
		{
			name:           "n=0",
			n:              0,
			randomSequence: nil,
			expected:       []Color{},
		},
		{
			name: "n=1",
			n:    1,
			// R1,G1,B1
			randomSequence: []int{100, 150, 200},
			expected: []Color{
				{R: 100, G: 150, B: 200},
			},
		},
		{
			name: "n=2",
			n:    2,
			// R1,G1,B1, R2,G2,B2
			randomSequence: []int{64, 65, 66, 254, 200, 128},
			expected: []Color{
				{R: 64, G: 65, B: 66},
				{R: 254, G: 200, B: 128},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mocks.NewMockRandom(ctrl)

			if test.n == 0 {
				emptyColors := InitColorsArray(test.n, r)
				assert.Equal(t, test.expected, emptyColors)
				return
			}

			// Ожидаем ровно 3*n вызовов GetRandomIntInterval(64,255),
			// и возвращаем заранее заданную последовательность значений.
			calls := 3 * test.n
			assert.Equal(t, calls, len(test.randomSequence), "randomSequence length must be 3*n")

			sequence := test.randomSequence
			idx := 0

			r.EXPECT().
				GetRandomIntInterval(64, 255).
				Times(calls).
				DoAndReturn(func(from, to int) int {
					v := sequence[idx]
					idx++
					return v
				})

			got := InitColorsArray(test.n, r)
			assert.Equal(t, test.expected, got)
		})
	}
}
