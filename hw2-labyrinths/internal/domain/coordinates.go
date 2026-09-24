package domain

type Coordinates struct {
	Row int
	Col int
}

func NewCoordinates(row, col int) Coordinates {
	return Coordinates{Row: row, Col: col}
}
