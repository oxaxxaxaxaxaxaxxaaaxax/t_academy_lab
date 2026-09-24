package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
)

func TestDrawMazeConsole(t *testing.T) {
	//enable unicode walls - false
	drawer := NewConsoleDrawer(false)

	type TestCase struct {
		pathToGraph string
		gridHeight  int
		gridWidth   int
	}

	testCases := []TestCase{
		{
			pathToGraph: "../../graphs/graphs.txt",
			gridHeight:  3,
			gridWidth:   3,
		},
		{
			pathToGraph: "../../graphs/graphs1.txt",
			gridHeight:  21,
			gridWidth:   21,
		},
		{
			pathToGraph: "../../graphs/graphs2.txt",
			gridHeight:  5,
			gridWidth:   5,
		},
	}

	for _, testCase := range testCases {
		graph, err := adapter.LoadGraphFromJSON(testCase.pathToGraph)
		assert.Nil(t, err)
		grid := drawer.DrawMaze(graph)
		assert.Equal(t, len(grid), testCase.gridHeight)
		assert.Equal(t, len(grid[0]), testCase.gridWidth)
	}

}
