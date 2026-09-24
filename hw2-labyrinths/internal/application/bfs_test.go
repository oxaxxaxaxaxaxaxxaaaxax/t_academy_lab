package application

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

func TestBFS(t *testing.T) {
	solver := NewBFSSolver()

	type TestCase struct {
		pathToGraph string
		parents     []int
	}

	testCases := []TestCase{
		{
			pathToGraph: "../../graphs/graphs.txt",
			parents:     []int{-1},
		},
		{
			pathToGraph: "../../graphs/graphs1.txt",
			parents: []int{-1, 0, 1, 2, 3, 4, 16, 17, 7, 19, 0, 10, 11, 23, 4, 5, 15, 27, 8, 29, 10, 20, 21, 22, 23,
				15, 25, 37, 38, 39, 31, 21, 31, 32, 44, 25, 46, 36, 37, 38, 30, 31, 41, 33, 43, 55, 45, 37, 49, 39, 40, 50, 42, 52, 44,
				54, 46, 56, 68, 58, 61, 62, 52, 53, 54, 75, 67, 57, 67, 68, 60, 72, 62, 83, 73, 74, 77, 67, 79, 69, 70, 71, 72, 82, 83, 84, 87, 77, 89,
				79, 80, 81, 82, 83, 84, 94, 86, 96, 88, 98},
		},
		{
			pathToGraph: "../../graphs/graphs2.txt",
			parents:     []int{-1, 3, 0, 2},
		},
	}

	for _, testCase := range testCases {
		graph, err := adapter.LoadGraphFromJSON(testCase.pathToGraph)
		assert.Nil(t, err, "Cannot load graph")
		parent, err := solver.SolveMaze(graph, domain.NewCoordinates(0, 0), domain.NewCoordinates(graph.Height, graph.Width))
		assert.Nil(t, err, "Cannot solve Maze")
		assert.True(t, reflect.DeepEqual(parent, testCase.parents), "Parents slices are not equal")

	}
}
