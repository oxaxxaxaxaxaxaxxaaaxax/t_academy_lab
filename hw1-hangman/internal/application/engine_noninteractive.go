package application

import (
	"handman/internal/domain"
)

type TestEngine struct {
	ui UI
}

func NewTestEngine(ui UI) TestEngine {
	return TestEngine{ui: ui}
}

func (engine TestEngine) StartGame() {
	hiddenWord := engine.ui.ReadHiddenWord()
	matchResult, matchStatus := domain.MatchWords(hiddenWord, engine.ui.ReadUserAnswer())
	engine.ui.ShowMatchResultWord(matchResult, matchStatus)
}
