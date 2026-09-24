package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	maxUint64    = math.MaxUint64
	minUint      = 0
	minStringLen = 0
)

func TestInspectorRand_RandBoll(t *testing.T) {
	r := NewInspectorRand()

	for i := 0; i < 50; i++ {
		v := r.RandBoll()
		assert.True(t, v == true || v == false)
	}
}

func TestInspectorRand_RandInt64(t *testing.T) {
	r := NewInspectorRand()

	for i := 0; i < 50; i++ {
		v := r.RandInt64()
		assert.GreaterOrEqual(t, v, int64(-intMax))
		assert.Less(t, v, int64(intMax))
	}
}

func TestInspectorRand_RandIntPositive(t *testing.T) {
	r := NewInspectorRand()

	for i := 0; i < 50; i++ {
		v := r.RandIntPositive()
		assert.GreaterOrEqual(t, v, minUint)
		assert.Less(t, v, intMax)
	}
}

func TestInspectorRand_RandUint64(t *testing.T) {
	r := NewInspectorRand()

	for i := 0; i < 50; i++ {
		v := r.RandUint64()
		assert.GreaterOrEqual(t, v, uint64(minUint))
		assert.Less(t, v, uint64(maxUint64))
	}
}

func TestInspectorRand_RandFloat64(t *testing.T) {
	r := NewInspectorRand()

	for i := 0; i < 50; i++ {
		v := r.RandFloat64()
		assert.GreaterOrEqual(t, v, -floatMax)
		assert.LessOrEqual(t, v, floatMax)
	}
}

func TestInspectorRand_RandString(t *testing.T) {
	r := NewInspectorRand()

	for i := 0; i < 50; i++ {
		s := r.RandString()

		assert.GreaterOrEqual(t, len([]rune(s)), minStringLen)
		assert.Less(t, len([]rune(s)), maxStringLen)
	}
}
