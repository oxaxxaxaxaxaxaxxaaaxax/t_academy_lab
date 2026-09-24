package application

type Algorithm int

func NewGenerator(algName string) Generator {
	switch algName {
	case "prim":
		return NewPrimGenerator()
	case "dfs":
		return NewDFSGenerator()
	default:
		return NewPrimGenerator()
	}
}

func NewSolver(algName string) Solver {
	switch algName {
	case "dijkstra":
		return NewDijkstraSolver()
	case "bfs":
		return NewBFSSolver()
	case "a_star":
		return NewAStarSolver()
	case "bellman_ford":
		return NewBellmanFordSolver()
	default:
		return NewDijkstraSolver()
	}
}

func NewDrawer(name string, enableUnicode bool) (Drawer, error) {
	switch name {
	case "console":
		return NewConsoleDrawer(enableUnicode), nil
	default:
		drawer, err := NewFileDrawer(name, enableUnicode)
		if err != nil {
			return nil, err
		}
		return drawer, nil
	}
}
