package application

import (
	"fmt"
	"os"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

type ConsoleDrawer struct {
	enableUnicode bool
}

func NewConsoleDrawer(enableUnicode bool) ConsoleDrawer {
	return ConsoleDrawer{enableUnicode: enableUnicode}
}

type SymbolsDrawer struct {
	Pass      rune
	Wall      rune
	StartCell rune
	EndCell   rune
	Path      rune
	Sand      rune
	Coin      rune
}

const (
	PassWeight = 2
	SandWeight = 5
	CoinWeight = 1
)

func NewSymbolsDrawer() SymbolsDrawer {
	return SymbolsDrawer{Pass: ' ', Wall: '#', StartCell: 'O',
		EndCell: 'E', Path: '*', Sand: '░', Coin: '¤'}
}

func CreateGrid(dr SymbolsDrawer, drawHeight, drawWidth int) [][]rune {
	grid := make([][]rune, drawHeight)
	for h := range drawHeight {
		grid[h] = make([]rune, drawWidth)
		for w := range drawWidth {
			grid[h][w] = dr.Wall
		}
	}
	return grid
}

func (c ConsoleDrawer) CreateTopologyGrid(drawHeight, drawWidth int) [][]uint8 {
	grid := make([][]uint8, drawHeight)
	for h := range drawHeight {
		grid[h] = make([]uint8, drawWidth)
	}
	return grid
}

func (c ConsoleDrawer) ShowMaze(graph domain.Graph) {
	grid := c.DrawMaze(graph)
	grid = c.DrawObstacle(grid, graph)
	c.WriteMaze(grid)
}

func (c ConsoleDrawer) DrawMaze(graph domain.Graph) [][]rune {
	drawHeight := graph.Height*2 + 1
	drawWidth := graph.Width*2 + 1

	topologyGrid := c.CreateTopologyGrid(drawHeight, drawWidth)

	for hIdx := range graph.Height {
		for wIdx := range graph.Width {
			topologyGrid = c.RenderCell(topologyGrid, graph, hIdx, wIdx)
		}
	}

	if c.enableUnicode {
		return FillUnicodeGrid(topologyGrid)
	}
	return FillGrid(topologyGrid)
}

func (c ConsoleDrawer) RenderCell(grid [][]uint8, graph domain.Graph, hIdx int, wIdx int) [][]uint8 {

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

func (c ConsoleDrawer) DrawSolution(graph domain.Graph, path []int, start, goal domain.Coordinates) {
	grid := c.DrawMaze(graph)
	grid = c.DrawObstacle(grid, graph)
	grid = DrawPath(grid, start, goal, path, graph.Width)
	grid = DrawStartAndGoalCell(start, goal, grid)
	c.WriteMaze(grid)
}

func DecodeCoordinates(idx int, width int) domain.CoordPair {
	return domain.CoordPair{idx / width, idx % width}
}

func (c ConsoleDrawer) DrawObstacle(grid [][]rune, graph domain.Graph) [][]rune {
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

func DrawPath(grid [][]rune, start, goal domain.Coordinates, path []int, width int) [][]rune {
	dr := NewSymbolsDrawer()

	startSl := domain.CoordPair{start.Row, start.Col}
	goalSl := domain.CoordPair{goal.Row, goal.Col}

	startNode := EncodeCoordinates(startSl, width)
	v := EncodeCoordinates(goalSl, width)
	for v != startNode {
		dest := DecodeCoordinates(v, width)
		drawDestHIdx := dest[0]*2 + 1
		drawDestWIdx := dest[1]*2 + 1
		grid[drawDestHIdx][drawDestWIdx] = dr.Path

		source := DecodeCoordinates(path[v], width)
		drawSourceHIdx := source[0]*2 + 1
		drawSourceWIdx := source[1]*2 + 1
		grid[drawSourceHIdx][drawSourceWIdx] = dr.Path

		middleHIdx := (drawDestHIdx + drawSourceHIdx) / 2
		middleWIdx := (drawDestWIdx + drawSourceWIdx) / 2

		grid[middleHIdx][middleWIdx] = dr.Path

		v = path[v]
	}
	return grid
}

func DrawStartAndGoalCell(start, goal domain.Coordinates, grid [][]rune) [][]rune {
	dr := NewSymbolsDrawer()

	drawStartRow := start.Row*2 + 1
	drawStartCol := start.Col*2 + 1
	drawGoalRow := goal.Row*2 + 1
	drawGoalCol := goal.Col*2 + 1

	grid[drawStartRow][drawStartCol] = dr.StartCell
	grid[drawGoalRow][drawGoalCol] = dr.EndCell

	return grid
}

func FillGrid(topologyGrid [][]uint8) [][]rune {
	drawer := NewSymbolsDrawer()
	drawHeight := len(topologyGrid)
	drawWidth := len(topologyGrid[0])
	grid := CreateGrid(drawer, drawHeight, drawWidth)

	for hIdx, h := range topologyGrid {
		for wIdx, cell := range h {
			if cell == 0 {
				grid[hIdx][wIdx] = drawer.Pass
			}
		}
	}
	return grid
}

func FillUnicodeGrid(topologyGrid [][]uint8) [][]rune {
	drawer := NewSymbolsDrawer()
	drawHeight := len(topologyGrid)
	drawWidth := len(topologyGrid[0])
	grid := CreateGrid(drawer, drawHeight, drawWidth)

	for hIdx := range drawHeight {
		for wIdx := range drawWidth {
			grid[hIdx][wIdx] = domain.GetSymbol(topologyGrid[hIdx][wIdx])
		}
	}
	return grid
}

func (c ConsoleDrawer) WriteMaze(grid [][]rune) {
	var b strings.Builder

	for _, c := range grid {
		for _, cell := range c {
			b.WriteRune(cell)
		}
		b.WriteRune('\n')
	}
	fmt.Fprintf(os.Stdout, "%s", b.String())
}
