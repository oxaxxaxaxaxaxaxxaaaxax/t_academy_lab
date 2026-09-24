package domain

type Cell struct {
	Row       int
	Col       int
	Walls     Direction
	IsVisited bool
	Weight    int
}

type Direction uint8

const (
	Right Direction = 1 << 0
	Down  Direction = 1 << 1
	Left  Direction = 1 << 2
	Up    Direction = 1 << 3
)

func NewCell(row, col int, walls Direction, weight int) Cell {
	return Cell{Row: row, Col: col, Walls: walls, IsVisited: false, Weight: weight}
}
