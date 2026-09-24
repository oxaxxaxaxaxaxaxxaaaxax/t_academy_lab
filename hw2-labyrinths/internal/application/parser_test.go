package application

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

func TestCreateCoordinates(t *testing.T) {
	type TestCase struct {
		cell                string
		expectedCoordinates domain.Coordinates
		expectedErr         error
	}

	tests := []TestCase{
		{
			cell:                "5,5",
			expectedCoordinates: domain.NewCoordinates(5, 5),
			expectedErr:         nil,
		},
		{
			cell:                "10,5",
			expectedCoordinates: domain.NewCoordinates(10, 5),
			expectedErr:         nil,
		},
		{
			cell:                "5k,!",
			expectedCoordinates: domain.Coordinates{},
			expectedErr:         ErrInvalidConfig,
		},
	}

	for _, test := range tests {
		coord, err := CreateCoordinates(test.cell)
		assert.Equal(t, test.expectedCoordinates, coord)
		assert.ErrorIs(t, test.expectedErr, err)
	}
}

func TestParseCmdLinePositive(f *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	type TestCase struct {
		cmdLine       []string
		Generate      bool
		Solve         bool
		Generator     string
		Solver        string
		Output        string
		File          string
		EnableUnicode bool
		Height        int
		Width         int
		Start         domain.Coordinates
		Goal          domain.Coordinates
		expectedErr   error
	}

	tests := []TestCase{
		{
			cmdLine: []string{
				"cmd.testPositive",
				"-generate",
				"-generate_algorithm=dfs",
				"-height=5",
				"-width=7",
				"-start=1,2",
			},
			Generate:      true,
			Solve:         false,
			Generator:     "dfs",
			Solver:        "dijkstra",
			EnableUnicode: false,
			Height:        5,
			Width:         7,
			Start:         domain.NewCoordinates(1, 2),
			Goal:          domain.NewCoordinates(4, 6),
			Output:        "console",
			File:          "maze1.txt",
			expectedErr:   nil,
		},
		{
			cmdLine: []string{
				"cmd.testWithoutFlags",
			},
			Generate:      false,
			Solve:         false,
			Generator:     "prim",
			Solver:        "dijkstra",
			EnableUnicode: false,
			Height:        10,
			Width:         10,
			File:          "maze1.txt",
			Output:        "console",
			Start:         domain.NewCoordinates(0, 0),
			Goal:          domain.NewCoordinates(9, 9),
			expectedErr:   nil,
		},
	}

	for _, test := range tests {
		os.Args = test.cmdLine
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		conf, err := ParseCmdLine()
		assert.ErrorIs(f, test.expectedErr, err)
		assert.Equal(f, test.Generate, conf.Generate)
		assert.Equal(f, test.Solve, conf.Solve)
		assert.Equal(f, test.Generator, conf.Generator)
		assert.Equal(f, test.Solver, conf.Solver)
		assert.Equal(f, test.Output, conf.Output)
		assert.Equal(f, test.EnableUnicode, conf.EnableUnicode)
		assert.Equal(f, test.Height, conf.Height)
		assert.Equal(f, test.Width, conf.Width)
		assert.Equal(f, test.Start, conf.Start)
		assert.Equal(f, test.Goal, conf.Goal)
		assert.Equal(f, test.EnableUnicode, conf.EnableUnicode)

	}
}

func TestParseCmdLineNegative(f *testing.T) {

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cmdLine := []string{
		"cmd.testNegative",
		"-generate",
		"-start=a,b",
		"-enable_unicode_walls",
	}
	expectedErr := ErrInvalidConfig

	os.Args = cmdLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	_, err := ParseCmdLine()
	assert.ErrorIs(f, expectedErr, err)
}

func TestParsePropertiesFile(t *testing.T) {
	type TestCase struct {
		filePath    string
		expectedCfg MazeConfig
		expectedErr error
	}

	tests := []TestCase{
		{
			filePath: "maze_pr1.txt",
			expectedCfg: MazeConfig{
				Height:   10,
				Width:    10,
				Filename: "maze1.txt",
				Pass:     ' ',
				Coin:     rune(0),
				Sand:     rune(0),
				Start:    domain.NewCoordinates(0, 0),
				Goal:     domain.NewCoordinates(9, 9),
			},
			expectedErr: nil,
		},
		{
			filePath: "maze_pr_a_star.txt",
			expectedCfg: MazeConfig{
				Height:   10,
				Width:    10,
				Filename: "maze_test.txt",
				Pass:     ' ',
				Sand:     '░',
				Coin:     '¤',
				Start:    domain.NewCoordinates(0, 0),
				Goal:     domain.NewCoordinates(9, 9),
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		cfg, err := ParsePropertiesFile(test.filePath)
		assert.ErrorIs(t, test.expectedErr, err)
		assert.Equal(t, test.expectedCfg.Height, cfg.Height)
		assert.Equal(t, test.expectedCfg.Width, cfg.Width)
		assert.Equal(t, test.expectedCfg.Filename, cfg.Filename)
		assert.Equal(t, test.expectedCfg.Pass, cfg.Pass)
		assert.Equal(t, test.expectedCfg.Coin, cfg.Coin)
		assert.Equal(t, test.expectedCfg.Sand, cfg.Sand)
		assert.Equal(t, test.expectedCfg.Start, cfg.Start)
		assert.Equal(t, test.expectedCfg.Goal, cfg.Goal)
	}
}

func TestParsePropertiesFileError(t *testing.T) {
	filePath := "maze_pr_wrong_start.txt"

	_, err := ParsePropertiesFile(filePath)
	assert.ErrorIs(t, ErrInvalidConfig, err)

}

func TestGenerateGraph_Size(t *testing.T) {
	type TestCase struct {
		mazeConf    MazeConfig
		expectedErr error
	}

	tests := []TestCase{
		{
			mazeConf: MazeConfig{
				Height:   10,
				Width:    10,
				Filename: "maze1.txt",
				Pass:     ' ',
				Coin:     rune(0),
				Sand:     rune(0),
				Start:    domain.NewCoordinates(0, 0),
				Goal:     domain.NewCoordinates(9, 9),
			},
			expectedErr: nil,
		},
		{
			mazeConf: MazeConfig{
				Height:   10,
				Width:    10,
				Filename: "maze_test.txt",
				Pass:     ' ',
				Sand:     '░',
				Coin:     '¤',
				Start:    domain.NewCoordinates(0, 0),
				Goal:     domain.NewCoordinates(9, 9),
			},
			expectedErr: nil,
		},
		{
			mazeConf: MazeConfig{
				Height:   5,
				Width:    5,
				Filename: "no_such_maze.txt",
			},
			expectedErr: ErrFileNotFound,
		},
	}

	for _, test := range tests {
		graph, err := GenerateGraph(test.mazeConf)

		if test.expectedErr != nil {
			assert.ErrorIs(t, err, test.expectedErr)
			return
		}

		assert.Equal(t, test.mazeConf.Height, graph.Height)
		assert.Equal(t, test.mazeConf.Width, graph.Width)
		assert.Len(t, graph.Graph, test.mazeConf.Height)
		assert.Len(t, graph.Graph[0], test.mazeConf.Width)
	}
}
