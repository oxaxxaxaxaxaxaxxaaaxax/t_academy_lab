package application

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type FileDrawer struct {
	filePath      string
	enableUnicode bool
}

func NewFileDrawer(fileName string, enableUnicode bool) (FileDrawer, error) {
	filePath := propertiesFilePath
	exist, err := exists(filePath)
	if err != nil {
		return FileDrawer{}, err
	}
	if exist {
		filePath = path.Join(filePath, fileName)
		return FileDrawer{filePath: filePath, enableUnicode: enableUnicode}, nil
	}
	err = os.Mkdir(filePath, 0755)
	if err != nil {
		panic(err)
	}
	filePath = path.Join(filePath, fileName)
	return FileDrawer{filePath: filePath, enableUnicode: enableUnicode}, nil
}

func (f FileDrawer) CreateTopologyGrid(drawHeight, drawWidth int) [][]uint8 {
	grid := make([][]uint8, drawHeight)
	for h := range drawHeight {
		grid[h] = make([]uint8, drawWidth)
	}
	return grid
}

func (f FileDrawer) ShowMaze(graph domain.Graph) {
	grid := f.DrawMaze(graph)
	grid = f.DrawObstacle(grid, graph)
	err := f.WriteMaze(grid)
	if err != nil {
		adapter.ShowError(err)
	}
}

func (f FileDrawer) DrawMaze(graph domain.Graph) [][]rune {
	drawHeight := graph.Height*2 + 1
	drawWidth := graph.Width*2 + 1

	topologyGrid := f.CreateTopologyGrid(drawHeight, drawWidth)

	for hIdx := range graph.Height {
		for wIdx := range graph.Width {
			topologyGrid = f.RenderCell(topologyGrid, graph, hIdx, wIdx)
		}
	}

	if f.enableUnicode {
		return FillUnicodeGrid(topologyGrid)
	}
	return FillGrid(topologyGrid)
}

func (f FileDrawer) FillGrid(topologyGrid [][]uint8) [][]rune {
	drawer := NewSymbolsDrawer()
	grid := CreateGrid(drawer, len(topologyGrid), len(topologyGrid[0]))

	for hIdx, h := range topologyGrid {
		for wIdx, cell := range h {
			if cell == 0 {
				grid[hIdx][wIdx] = drawer.Pass
			}
		}
	}
	return grid
}

func (f FileDrawer) DrawSolution(graph domain.Graph, path []int, start, goal domain.Coordinates) {
	grid := f.DrawMaze(graph)
	grid = f.DrawObstacle(grid, graph)
	grid = DrawPath(grid, start, goal, path, graph.Width)
	grid = DrawStartAndGoalCell(start, goal, grid)
	err := f.WriteMaze(grid)
	if err != nil {
		adapter.ShowError(err)
	}
}

func (f FileDrawer) DrawObstacle(grid [][]rune, graph domain.Graph) [][]rune {
	dr := NewSymbolsDrawer()
	for hIdx := range graph.Height {
		for wIdx := range graph.Width {
			drawRow := hIdx*2 + 1
			drawCol := wIdx*2 + 1

			switch {
			case graph.Graph[hIdx][wIdx].Weight == SandWeight:
				grid[drawRow][drawCol] = dr.Sand
			case graph.Graph[hIdx][wIdx].Weight == CoinWeight:
				grid[drawRow][drawCol] = dr.Coin
			default:
				continue
			}
		}
	}
	return grid
}

func (f FileDrawer) WriteMaze(grid [][]rune) error {
	file, err := os.Create(f.filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			adapter.ShowError(err)
		}
	}()
	w := bufio.NewWriter(file)
	for _, c := range grid {
		for _, cell := range c {
			_, err = w.WriteRune(cell)
			if err != nil {
				panic(err)
			}
		}
		_, err = w.WriteRune('\n')
		if err != nil {
			panic(err)
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func (f FileDrawer) RenderCell(grid [][]uint8, graph domain.Graph, hIdx int, wIdx int) [][]uint8 {

	neighborWalls := map[string]uint8{"LU": 0, "U": 0, "RU": 0, "L": 0,
		"R": 0, "LD": 0, "D": 0, "RD": 0}

	cell := graph.Graph[hIdx][wIdx]

	if !InGrid(graph.Width, graph.Height, hIdx-1, wIdx) || cell.Walls&domain.Up != 0 {
		neighborWalls["U"] |= uint8(domain.Left) | uint8(domain.Right)
		neighborWalls["LU"] |= uint8(domain.Right)
		neighborWalls["RU"] |= uint8(domain.Left)
	}

	if !InGrid(graph.Width, graph.Height, hIdx, wIdx-1) || cell.Walls&domain.Left != 0 {
		neighborWalls["L"] |= uint8(domain.Up) | uint8(domain.Down)
		neighborWalls["LU"] |= uint8(domain.Down)
		neighborWalls["LD"] |= uint8(domain.Up)
	}

	if !InGrid(graph.Width, graph.Height, hIdx, wIdx+1) || cell.Walls&domain.Right != 0 {
		neighborWalls["R"] |= uint8(domain.Up) | uint8(domain.Down)
		neighborWalls["RU"] |= uint8(domain.Down)
		neighborWalls["RD"] |= uint8(domain.Up)
	}

	if !InGrid(graph.Width, graph.Height, hIdx+1, wIdx) || cell.Walls&domain.Down != 0 {
		neighborWalls["D"] |= uint8(domain.Left) | uint8(domain.Right)
		neighborWalls["LD"] |= uint8(domain.Right)
		neighborWalls["RD"] |= uint8(domain.Left)
	}

	drawHIdx := hIdx*2 + 1
	drawWIdx := wIdx*2 + 1

	grid[drawHIdx][drawWIdx] = 0
	grid[drawHIdx-1][drawWIdx-1] |= neighborWalls["LU"]
	grid[drawHIdx-1][drawWIdx] |= neighborWalls["U"]
	grid[drawHIdx-1][drawWIdx+1] |= neighborWalls["RU"]
	grid[drawHIdx][drawWIdx-1] |= neighborWalls["L"]
	grid[drawHIdx][drawWIdx+1] |= neighborWalls["R"]
	grid[drawHIdx+1][drawWIdx-1] |= neighborWalls["LD"]
	grid[drawHIdx+1][drawWIdx] |= neighborWalls["D"]
	grid[drawHIdx+1][drawWIdx+1] |= neighborWalls["RD"]

	return grid
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
