package application

import (
	"handman/internal/domain"
)

type GameMode int

const (
	Test GameMode = iota
	Interactive
)

func NewGame(ui UI, rnd Rand) domain.Engine {
	lengthArgs := ui.GetLengthInput()

	if lengthArgs == 3 {
		return getGameByMode(Test, ui, rnd)
	} else {
		return getGameByMode(Interactive, ui, rnd)
	}
}

func getGameByMode(mode GameMode, ui UI, rnd Rand) domain.Engine {
	switch mode {
	case Interactive:
		return NewInteractiveEngine(ui, rnd)
	case Test:
		return NewTestEngine(ui)
	default:
		return NewInteractiveEngine(ui, rnd)
	}
}
