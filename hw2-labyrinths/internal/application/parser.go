package application

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

const propertiesFilePath = "../../maze_properties"

var (
	ErrFileNotFound  = errors.New("file not found")
	ErrInvalidConfig = errors.New("invalid config")
)

type Config struct {
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
}

func ParseCmdLine() (Config, error) {
	generate := flag.Bool("generate", false, "generate mode")
	solve := flag.Bool("solve", false, "solve mode")
	generatorAlgorithm := flag.String("generate_algorithm", "prim", "algorithm to generate maze")
	solverAlgorithm := flag.String("solve_algorithm", "dijkstra", "algorithm to solve maze")
	file := flag.String("file", "maze1.txt", "file contains the maze to solve")
	out := flag.String("output", "console", "file to write maze")
	unicode := flag.Bool("enable_unicode_walls", false, "enable maze with unicode symbols walls")
	h := flag.Int("height", 10, "height of maze")
	w := flag.Int("width", 10, "width of maze")
	start := flag.String("start", "0,0", "start position")
	goal := flag.String("goal", "", "goal position")

	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, `
		Help message for usage:

		-generate 
			Start generate mode 
		-solve
			Start solve mode
		-generate_algorithm (string)
			Algorithm used to generate the maze. Supported: prim, dfs (default "prim")
		-solve_algorithm (string)
			Algorithm used to solve the maze. Supported: dijkstra, astar (default "dijkstra")
		-file (string)
			File contains maze description for solving
		-output (string)
			Output target for the maze: "console" or a file name (default "console")
		-enable_unicode_walls
			Render maze walls using Unicode symbols. Disabled by default.
		-height (int)
			Maze height in cells (default 10)
		-width (int)
			Maze width in cells (default 10)
		-start (string)
			Starting cell coordinates in the format "row,col" (default "0,0")
		-goal string
			Goal cell coordinates in the format "row,col" (default "height-1,width-1")`)
		if err != nil {
			adapter.ShowError(err)
		}
	}

	flag.Parse()

	startCoord, err := CreateCoordinates(*start)
	if err != nil {
		return Config{}, err
	}

	//if user don't enter goal coordinates
	if *goal == "" {
		*goal = strconv.Itoa((*h)-1) + "," + strconv.Itoa((*w)-1)
	}
	goalCoord, err := CreateCoordinates(*goal)
	if err != nil {
		return Config{}, err
	}

	return Config{Generate: *generate, Solve: *solve,
		Generator: *generatorAlgorithm, Solver: *solverAlgorithm,
		Output: *out, File: *file, EnableUnicode: *unicode, Height: *h,
		Width: *w, Start: startCoord, Goal: goalCoord}, nil
}

type MazeConfig struct {
	Height   int
	Width    int
	Filename string
	Pass     rune
	Sand     rune
	Coin     rune
	Start    domain.Coordinates
	Goal     domain.Coordinates
}

func ParseMaze(fileName string) (domain.Graph, MazeConfig, error) {
	mazeConf, err := ParsePropertiesFile(fileName)
	if err != nil {
		return domain.Graph{}, MazeConfig{}, err
	}
	graph, err := GenerateGraph(mazeConf)
	if err != nil {
		return domain.Graph{}, MazeConfig{}, err
	}
	return graph, mazeConf, nil
}

func ParsePropertiesFile(fileName string) (MazeConfig, error) {
	filePath := propertiesFilePath
	filePath = path.Join(filePath, fileName)
	if exists, _ := exists(filePath); exists == false {
		return MazeConfig{}, ErrFileNotFound
	}
	file, err := os.Open(filePath)
	if err != nil {
		return MazeConfig{}, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			adapter.ShowError(err)
		}
	}()
	var height int
	var width int
	var mazeFilename string
	var pass rune
	var coin rune
	var sand rune
	var start, goal domain.Coordinates

	r := bufio.NewScanner(file)

	for r.Scan() {
		line := r.Text()
		kv := strings.Split(line, "=")
		for idx := range kv {
			kv[idx] = strings.TrimSpace(kv[idx])
		}
		if len(kv) != 2 {
			return MazeConfig{}, ErrInvalidConfig
		}
		key, value := kv[0], kv[1]
		switch key {
		case "width":
			width, err = strconv.Atoi(value)
			if err != nil {
				return MazeConfig{}, ErrInvalidConfig
			}
		case "height":
			height, err = strconv.Atoi(value)
			if err != nil {
				return MazeConfig{}, errors.Join(ErrInvalidConfig, err)
			}
		case "pass":
			str, err := strconv.Unquote(value)
			if err != nil {
				return MazeConfig{}, errors.Join(ErrInvalidConfig, err)
			}
			pass = []rune(str)[0]
		case "coin":
			coin = []rune(value)[0]
		case "sand":
			sand = []rune(value)[0]
		case "maze_filename":
			mazeFilename = value
		case "start":
			start, err = CreateCoordinates(value)
			if err != nil || start.Row >= height || start.Col >= width {
				return MazeConfig{}, ErrInvalidConfig
			}
		case "goal":
			goal, err = CreateCoordinates(value)
			if err != nil || goal.Row >= height || goal.Col >= width {
				return MazeConfig{}, ErrInvalidConfig
			}
		}
	}
	return MazeConfig{Height: height, Width: width, Filename: mazeFilename,
		Pass: pass, Sand: sand, Coin: coin, Start: start, Goal: goal}, nil
}

func GenerateGraph(mazeConf MazeConfig) (domain.Graph, error) {
	height := mazeConf.Height
	width := mazeConf.Width
	pass := mazeConf.Pass
	filename := mazeConf.Filename
	dr := NewSymbolsDrawer()

	filePath := propertiesFilePath
	filePath = path.Join(filePath, filename)
	if exists, _ := exists(filePath); exists == false {
		return domain.Graph{}, ErrFileNotFound
	}
	file, err := os.Open(filePath)
	if err != nil {
		return domain.Graph{}, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			adapter.ShowError(errors.Join(err))
		}
	}()
	graphMaze := CreateGraph(height, width)
	graph := graphMaze.Graph

	grid := make([][]rune, 0)

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		grid = append(grid, []rune(line))
	}

	for hIdx := range height {
		for wIdx := range width {
			gridHIdx := hIdx*2 + 1
			gridWIdx := wIdx*2 + 1

			cell := grid[gridHIdx][gridWIdx]
			switch {
			case cell == dr.Sand:
				graph[hIdx][wIdx].Weight = SandWeight
			case cell == dr.Coin:
				graph[hIdx][wIdx].Weight = CoinWeight
			default:
				graph[hIdx][wIdx].Weight = PassWeight
			}

			upNeighbors := grid[gridHIdx-1][gridWIdx]
			downNeighbors := grid[gridHIdx+1][gridWIdx]
			leftNeighbors := grid[gridHIdx][gridWIdx-1]
			rightNeighbors := grid[gridHIdx][gridWIdx+1]

			if upNeighbors == pass {
				graph[hIdx][wIdx].Walls &= ^domain.Up
			}
			if downNeighbors == pass {
				graph[hIdx][wIdx].Walls &= ^domain.Down
			}
			if leftNeighbors == pass {
				graph[hIdx][wIdx].Walls &= ^domain.Left
			}
			if rightNeighbors == pass {
				graph[hIdx][wIdx].Walls &= ^domain.Right
			}
			graph[hIdx][wIdx].Row = hIdx
			graph[hIdx][wIdx].Col = wIdx
		}
	}
	return graphMaze, nil

}

func CreateCoordinates(cell string) (domain.Coordinates, error) {
	cords := strings.Split(cell, ",")
	row, err := strconv.Atoi(cords[0])
	if err != nil {
		return domain.Coordinates{}, ErrInvalidConfig
	}
	col, err := strconv.Atoi(cords[1])
	if err != nil {
		return domain.Coordinates{}, ErrInvalidConfig
	}
	return domain.NewCoordinates(row, col), nil
}
