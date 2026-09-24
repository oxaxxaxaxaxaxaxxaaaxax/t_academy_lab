package application

import (
	"handman/internal/domain"
	"handman/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestInteractiveEngineCorrectAttemptNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ShowAttemptSettings().Times(1),
		ui.EXPECT().ReadUserNumber().Return(6, nil).Times(1))

	rnd := mocks.NewMockRand(ctrl)
	e := NewInteractiveEngine(ui, rnd)
	att, err := e.SelectAttemptCount()
	if att != 6 || err != nil {
		t.Fatalf("Failed to select attempt count: got %v, %v, want %v, %v", att, err, 6, nil)
	}
}

func TestInteractiveEngineIncorrectAttemptNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().ShowAttemptSettings().Times(1),
		ui.EXPECT().ReadUserNumber().Return(3, nil).Times(1))

	rnd := mocks.NewMockRand(ctrl)
	e := NewInteractiveEngine(ui, rnd)
	att, err := e.SelectAttemptCount()
	if err == nil {
		t.Fatalf("Expected error for bad input: got %v, %v, want %v, %v", att, err, 6, nil)
	}
}

func TestInteractiveEngineCorrectCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().SelectWordCategory().Times(1),
		ui.EXPECT().ReadUserRune().Return('1', nil).Times(1))

	rnd := mocks.NewMockRand(ctrl)
	e := NewInteractiveEngine(ui, rnd)

	if cat, err := e.SelectWordCategory(); cat != 1 || err != nil {
		t.Fatalf("Failed to select word category: got %v, %v, want %v, %v", cat, err, 1, nil)
	}
}

func TestInteractiveEngineCorrectDifficulty(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)

	gomock.InOrder(
		ui.EXPECT().SelectDifficultyLevel().Times(1),
		ui.EXPECT().ReadUserRune().Return('1', nil).Times(1))

	rnd := mocks.NewMockRand(ctrl)
	e := NewInteractiveEngine(ui, rnd)

	if cat, err := e.SelectLevelDifficulty(); cat != 1 || err != nil {
		t.Fatalf("Failed to select level difficulty: got %v, %v, want %v, %v", cat, err, 1, nil)
	}
}

func TestInteractiveEngineStartGameUpdateAttr(t *testing.T) {
	type ts struct {
		attr              GameAttr
		matchResult       string
		stat              domain.Status
		expectedReturnVal ReturnStatus
		callOrder         func(ui *mocks.MockUI)
	}

	cases := []ts{
		{
			attr:              GameAttr{UserWord: "sp*ce", HiddenWord: "space", CurrentImage: Default, WordLength: 5, CountGuessLetter: 4, CountAttempts: 6, CountMistakes: 0},
			matchResult:       "space",
			stat:              domain.Success,
			expectedReturnVal: Win,
			callOrder: func(ui *mocks.MockUI) {
				ui.EXPECT().ShowWord("space").Times(1)
			},
		},
		{
			attr:              GameAttr{UserWord: "sp*ce", HiddenWord: "space", CurrentImage: RightLeg, WordLength: 5, CountGuessLetter: 4, CountAttempts: 6, CountMistakes: 5},
			matchResult:       "sp*ce",
			stat:              domain.Fail,
			expectedReturnVal: Lose,
			callOrder: func(ui *mocks.MockUI) {
				ui.EXPECT().ShowWord("sp*ce").Times(1)
				ui.EXPECT().ShowCountMistakes(6, 6).Times(1)
				ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1)
			},
		},
		{
			attr:              GameAttr{UserWord: "sp*ce", HiddenWord: "space", CurrentImage: RightHand, WordLength: 5, CountGuessLetter: 4, CountAttempts: 6, CountMistakes: 3},
			matchResult:       "sp*ce",
			stat:              domain.Fail,
			expectedReturnVal: Continue,
			callOrder: func(ui *mocks.MockUI) {
				ui.EXPECT().ShowDescription("the universe")
				ui.EXPECT().ShowWord("sp*ce").Times(1)
				ui.EXPECT().ShowCountMistakes(6, 4).Times(1)
				ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1)
			},
		},
	}

	for _, val := range cases {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		ui := mocks.NewMockUI(ctrl)
		rnd := mocks.NewMockRand(ctrl)
		val.callOrder(ui)
		e := NewInteractiveEngine(ui, rnd)
		if st := e.UpdateAttr(&val.attr, val.matchResult, val.stat); st != val.expectedReturnVal {
			t.Errorf("Failed Return status: got %v, want %v", st, val.expectedReturnVal)
		}
	}
}

func TestInteractiveEngineStartGameInitializeAttr(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	ui := mocks.NewMockUI(ctrl)
	rnd := mocks.NewMockRand(ctrl)
	e := NewInteractiveEngine(ui, rnd)
	testAttr := e.InitializeAttr([]rune("saturn"), 6)

	if testAttr.WordLength != 6 || testAttr.UserWord != "******" ||
		testAttr.CountMistakes != 0 || testAttr.CountAttempts != 6 || testAttr.CountGuessLetter != 0 ||
		testAttr.CurrentImage != Default || testAttr.HiddenWord != "saturn" {
		t.Fatalf("Failed initialize game attribute: got %v,%v,%v,%v,%v,%v,%v want %v,%v,%v,%v,%v,%v,%v", testAttr.HiddenWord,
			testAttr.UserWord, testAttr.CountMistakes, testAttr.CountAttempts, testAttr.CountGuessLetter,
			testAttr.CurrentImage, testAttr.WordLength, "saturn", "******", 0, 6, 0, Default, 6)
	}
}

func TestInteractiveEngineStartGame(t *testing.T) {
	type ts struct {
		name      string
		callOrder func(ui *mocks.MockUI, rnd *mocks.MockRand)
	}

	cases := []ts{
		{
			name: "Win case",
			callOrder: func(ui *mocks.MockUI, rnd *mocks.MockRand) {
				gomock.InOrder(
					ui.EXPECT().ShowAttemptSettings().Times(1),
					ui.EXPECT().ReadUserNumber().Return(6, nil).Times(1),
					ui.EXPECT().SelectWordCategory().Times(1),
					ui.EXPECT().ReadUserRune().Return('0', nil).Times(1),
					ui.EXPECT().SelectDifficultyLevel().Times(1),
					ui.EXPECT().ReadUserRune().Return('0', nil).Times(1),
					ui.EXPECT().ShowWordCategory(gomock.Any()).Times(1),
					ui.EXPECT().ShowDifficultyLevel(gomock.Any()).Times(1),
					rnd.EXPECT().GetRandomWord(0, 0).Return([]rune("cat"), nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ReadUserRune().Return('c', nil).Times(1),
					ui.EXPECT().ShowWord("c**").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('a', nil).Times(1),
					ui.EXPECT().ShowWord("ca*").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('t', nil).Times(1),
					ui.EXPECT().ShowWord("cat").Times(1),
					ui.EXPECT().ShowWin().Times(1))
			},
		},
		{
			name: "Random category case",
			callOrder: func(ui *mocks.MockUI, rnd *mocks.MockRand) {
				gomock.InOrder(
					ui.EXPECT().ShowAttemptSettings().Times(1),
					ui.EXPECT().ReadUserNumber().Return(6, nil).Times(1),
					ui.EXPECT().SelectWordCategory().Times(1),
					ui.EXPECT().ReadUserRune().Return('5', nil).Times(1),
					ui.EXPECT().SelectDifficultyLevel().Times(1),
					ui.EXPECT().ReadUserRune().Return('0', nil).Times(1),
					ui.EXPECT().ShowErrorMessage(gomock.Any(), ErrBadCategory).Times(1),
					rnd.EXPECT().GetRandomCategory().Times(1),
					ui.EXPECT().ShowWordCategory(gomock.Any()).Times(1),
					ui.EXPECT().ShowDifficultyLevel(gomock.Any()).Times(1),
					rnd.EXPECT().GetRandomWord(0, 0).Return([]rune("cat"), nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ReadUserRune().Return('c', nil).Times(1),
					ui.EXPECT().ShowWord("c**").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('a', nil).Times(1),
					ui.EXPECT().ShowWord("ca*").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('t', nil).Times(1),
					ui.EXPECT().ShowWord("cat").Times(1),
					ui.EXPECT().ShowWin().Times(1))
			},
		},
		{
			name: "Random difficulty case",
			callOrder: func(ui *mocks.MockUI, rnd *mocks.MockRand) {
				gomock.InOrder(
					ui.EXPECT().ShowAttemptSettings().Times(1),
					ui.EXPECT().ReadUserNumber().Return(6, nil).Times(1),
					ui.EXPECT().SelectWordCategory().Times(1),
					ui.EXPECT().ReadUserRune().Return('0', nil).Times(1),
					ui.EXPECT().SelectDifficultyLevel().Times(1),
					ui.EXPECT().ReadUserRune().Return('9', nil).Times(1),
					ui.EXPECT().ShowErrorMessage(gomock.Any(), ErrBadDifficulty).Times(1),
					rnd.EXPECT().GetRandomDifficulty().Times(1),
					ui.EXPECT().ShowWordCategory(gomock.Any()).Times(1),
					ui.EXPECT().ShowDifficultyLevel(gomock.Any()).Times(1),
					rnd.EXPECT().GetRandomWord(0, 0).Return([]rune("cat"), nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ReadUserRune().Return('c', nil).Times(1),
					ui.EXPECT().ShowWord("c**").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('a', nil).Times(1),
					ui.EXPECT().ShowWord("ca*").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('t', nil).Times(1),
					ui.EXPECT().ShowWord("cat").Times(1),
					ui.EXPECT().ShowWin().Times(1))
			},
		},
		{
			name: "Random category and difficulty case",
			callOrder: func(ui *mocks.MockUI, rnd *mocks.MockRand) {
				gomock.InOrder(
					ui.EXPECT().ShowAttemptSettings().Times(1),
					ui.EXPECT().ReadUserNumber().Return(6, nil).Times(1),
					ui.EXPECT().SelectWordCategory().Times(1),
					ui.EXPECT().ReadUserRune().Return('9', nil).Times(1),
					ui.EXPECT().SelectDifficultyLevel().Times(1),
					ui.EXPECT().ReadUserRune().Return('9', nil).Times(1),
					ui.EXPECT().ShowErrorMessage(gomock.Any(), ErrBadDifficultyAndCategory).Times(1),
					rnd.EXPECT().GetRandomCategory().Times(1),
					rnd.EXPECT().GetRandomDifficulty().Times(1),
					ui.EXPECT().ShowWordCategory(gomock.Any()).Times(1),
					ui.EXPECT().ShowDifficultyLevel(gomock.Any()).Times(1),
					rnd.EXPECT().GetRandomWord(0, 0).Return([]rune("cat"), nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ReadUserRune().Return('c', nil).Times(1),
					ui.EXPECT().ShowWord("c**").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('a', nil).Times(1),
					ui.EXPECT().ShowWord("ca*").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 0).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('t', nil).Times(1),
					ui.EXPECT().ShowWord("cat").Times(1),
					ui.EXPECT().ShowWin().Times(1))
			},
		},
		{
			name: "Bad attempt settings case",
			callOrder: func(ui *mocks.MockUI, rnd *mocks.MockRand) {
				gomock.InOrder(
					ui.EXPECT().ShowAttemptSettings().Times(1),
					ui.EXPECT().ReadUserNumber().Return(3, nil).Times(1),
					ui.EXPECT().ShowBadInput().Times(1),
					ui.EXPECT().ShowError(gomock.Any()).Times(1))
			},
		},
		{
			name: "Lose case",
			callOrder: func(ui *mocks.MockUI, rnd *mocks.MockRand) {
				gomock.InOrder(
					ui.EXPECT().ShowAttemptSettings().Times(1),
					ui.EXPECT().ReadUserNumber().Return(6, nil).Times(1),
					ui.EXPECT().SelectWordCategory().Times(1),
					ui.EXPECT().ReadUserRune().Return('0', nil).Times(1),
					ui.EXPECT().SelectDifficultyLevel().Times(1),
					ui.EXPECT().ReadUserRune().Return('0', nil).Times(1),
					ui.EXPECT().ShowWordCategory(gomock.Any()).Times(1),
					ui.EXPECT().ShowDifficultyLevel(gomock.Any()).Times(1),
					rnd.EXPECT().GetRandomWord(0, 0).Return([]rune("cat"), nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ReadUserRune().Return('g', nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 1).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('p', nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 2).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('v', nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 3).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('e', nil).Times(1),
					ui.EXPECT().ShowDescription("^._.^").Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 4).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('q', nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 5).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ReadUserRune().Return('z', nil).Times(1),
					ui.EXPECT().ShowWord("***").Times(1),
					ui.EXPECT().ShowCountMistakes(6, 6).Times(1),
					ui.EXPECT().ShowHangmanPicture(gomock.Any()).Times(1),
					ui.EXPECT().ShowLose().Times(1))
			},
		},
	}

	for _, val := range cases {
		t.Run(val.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			ui := mocks.NewMockUI(ctrl)
			rnd := mocks.NewMockRand(ctrl)
			val.callOrder(ui, rnd)
			e := NewInteractiveEngine(ui, rnd)
			e.StartGame()
		})
	}
}
