package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDFS(t *testing.T) {
	generator := NewDFSGenerator()

	type TestCase struct {
		height, width int
	}

	testCases := []TestCase{
		{height: 2, width: 2},
		{height: 10, width: 10},
		{height: 1, width: 1},
	}

	for _, testCase := range testCases {
		graph, err := generator.GenerateMaze(testCase.height, testCase.width)
		assert.Nil(t, err, "Cannot generate Maze")
		assert.Equal(t, graph.Height, testCase.height, "Height maze is not equal")
		assert.Equal(t, graph.Width, testCase.width, "Width maze is not equal")
	}
}
