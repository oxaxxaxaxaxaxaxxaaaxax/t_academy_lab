package domain

var Symbols []rune = []rune{' ', '╶', '╷', '┌', '╴', '─', '┐', '┬', '╵',
	'└', '│', '├', '┘', '┴', '┤', '┼'}

func GetSymbol(idx uint8) rune {
	return Symbols[idx]
}
