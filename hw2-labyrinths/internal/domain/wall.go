package domain

type Wall struct {
	CellFromRow int
	CellFromCol int
	CellToRow   int
	CellToCol   int
}

func NewWall(From, To Cell) Wall {
	return Wall{CellFromRow: From.Row, CellToRow: To.Row,
		CellFromCol: From.Col, CellToCol: To.Col}
}
