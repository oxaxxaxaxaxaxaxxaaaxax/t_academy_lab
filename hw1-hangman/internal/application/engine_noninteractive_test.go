package application

import (
	"handman/internal/domain"
	"handman/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestTestEngineStartGameNeg(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ReadHiddenWord().Return("волокно").Times(1),
		ui.EXPECT().ReadUserAnswer().Return("толокно").Times(1),
		ui.EXPECT().ShowMatchResultWord("*олокно", domain.Fail).Times(1))

	e := NewTestEngine(ui)
	e.StartGame()
}

func TestTestEngineStartGamePos(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ReadHiddenWord().Return("cat").Times(1),
		ui.EXPECT().ReadUserAnswer().Return("cat").Times(1),
		ui.EXPECT().ShowMatchResultWord("cat", domain.Success).Times(1))
	e := NewTestEngine(ui)
	e.StartGame()
}

func TestTestEngineStartGameEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ReadHiddenWord().Return("").Times(1),
		ui.EXPECT().ReadUserAnswer().Return("").Times(1),
		ui.EXPECT().ShowMatchResultWord("", domain.Success).Times(1))
	e := NewTestEngine(ui)
	e.StartGame()
}

func TestTestEngineStartGameNormalize(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ReadHiddenWord().Return("cat").Times(1),
		ui.EXPECT().ReadUserAnswer().Return("Cat").Times(1),
		ui.EXPECT().ShowMatchResultWord("cat", domain.Success).Times(1))
	e := NewTestEngine(ui)
	e.StartGame()
}

func TestTestEngineStartGameFullLetters(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ReadHiddenWord().Return("кулон").Times(1),
		ui.EXPECT().ReadUserAnswer().Return("клоун").Times(1),
		ui.EXPECT().ShowMatchResultWord("кулон", domain.Success).Times(1))
	e := NewTestEngine(ui)
	e.StartGame()
}
