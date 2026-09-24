package domain

type CoordPair = [2]int
type Edge struct {
	NodeSource [2]int
	NodeDest   [2]int
	Weight     int
}

func NewEdge(rSource, cSource, rDest, cDest, weight int) Edge {
	return Edge{NodeSource: [2]int{rSource, cSource},
		NodeDest: [2]int{rDest, cDest}, Weight: weight}
}

func (e Edge) IsEqual(b Edge) bool {
	return (e.NodeSource == b.NodeSource && e.NodeDest == b.NodeDest) ||
		(e.NodeSource == b.NodeDest && e.NodeDest == b.NodeSource)
}
