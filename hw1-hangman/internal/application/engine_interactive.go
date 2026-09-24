package application

import (
	"errors"
	"handman/internal/adapter"
	"handman/internal/domain"
)

type ReturnStatus int

const (
	Continue ReturnStatus = iota
	Win
	Lose
)

type HangmanImage int

const (
	Default HangmanImage = iota
	Head
	Body
	RightHand
	LeftHand
	RightLeg
	LeftLeg
	HangmanImageCount
)

type GameAttr struct {
	UserWord         string
	HiddenWord       string
	CurrentImage     HangmanImage
	WordLength       int
	CountGuessLetter int
	CountAttempts    int
	CountMistakes    int
}

var (
	ErrBadCategory              = errors.New("Bad category")
	ErrBadDifficulty            = errors.New("Bad difficulty")
	ErrBadDifficultyAndCategory = errors.New("Bad difficulty and category")
)

type InteractiveEngine struct {
	ui  UI
	rnd Rand
}

func NewInteractiveEngine(ui UI, rnd Rand) InteractiveEngine {
	return InteractiveEngine{ui: ui, rnd: rnd}
}

func (engine InteractiveEngine) StartGame() {
	var returnStat ReturnStatus
	attCount, err := engine.SelectAttemptCount()
	if err != nil {
		engine.ui.ShowBadInput()
		engine.ui.ShowError(err)
		return
	}

	hiddenWord := engine.GetRandomWordByParams()
	gameAttr := engine.InitializeAttr(hiddenWord, attCount)
	engine.ui.ShowWord(gameAttr.UserWord)

	for gameAttr.CountMistakes != gameAttr.CountAttempts {
		letter, err := engine.ui.ReadUserRune()
		if err != nil {
			engine.ui.ShowBadInput()
			return
		}
		matchResult, matchStatus := domain.MatchLetter(letter, gameAttr.UserWord, string(hiddenWord))
		returnStat = engine.UpdateAttr(gameAttr, matchResult, matchStatus)

		switch returnStat {
		case Win:
			engine.ui.ShowWin()
			return
		case Lose:
			engine.ui.ShowLose()
		default:
			continue
		}
	}
}

func (InteractiveEngine) InitializeAttr(hiddenWord []rune, attCount int) *GameAttr {
	initialWord := make([]rune, len(hiddenWord))

	for i := range len(hiddenWord) {
		initialWord[i] = '*'
	}

	attr := GameAttr{
		HiddenWord:       string(hiddenWord),
		UserWord:         string(initialWord),
		CurrentImage:     Default,
		WordLength:       len(hiddenWord),
		CountGuessLetter: 0,
		CountAttempts:    attCount,
		CountMistakes:    0,
	}
	return &attr
}

func (engine InteractiveEngine) UpdateAttr(attr *GameAttr, guessWord string, stat domain.Status) ReturnStatus {
	attr.CountGuessLetter += domain.CountDiffSymbols([]rune(guessWord), []rune(attr.UserWord))
	attr.UserWord = guessWord

	if attr.CountGuessLetter == attr.WordLength {
		engine.ui.ShowWord(attr.UserWord)
		return Win
	}

	if stat == domain.Fail {
		attr.CountMistakes++
		attr.CurrentImage = HangmanImage(domain.CurrentPictureIndex(attr.CountMistakes, int(HangmanImageCount), attr.CountAttempts))

		if attr.CountMistakes == attr.CountAttempts {
			engine.ui.ShowWord(attr.UserWord)
			engine.ui.ShowCountMistakes(attr.CountAttempts, attr.CountMistakes)
			engine.ui.ShowHangmanPicture(adapter.GetHangmanPicture(int(attr.CurrentImage)))
			return Lose
		} else {
			attr.CurrentImage = HangmanImage(domain.CurrentPictureIndex(attr.CountMistakes, int(HangmanImageCount), attr.CountAttempts))
		}

		if attr.CountMistakes == (attr.CountAttempts*2)/3 {
			desc := domain.GetHelpDescription(attr.HiddenWord)
			engine.ui.ShowDescription(desc)
		}
	}

	engine.ui.ShowWord(attr.UserWord)
	engine.ui.ShowCountMistakes(attr.CountAttempts, attr.CountMistakes)
	engine.ui.ShowHangmanPicture(adapter.GetHangmanPicture(int(attr.CurrentImage)))
	return Continue
}

func (engine InteractiveEngine) SelectCategoryAndDifficulty() (int, int, error) {
	cat, errCat := engine.SelectWordCategory()
	lev, errDiff := engine.SelectLevelDifficulty()

	if errDiff != nil && errCat != nil {
		return 0, 0, ErrBadDifficultyAndCategory
	}
	if errCat != nil {
		return 0, lev, errCat
	}
	if errDiff != nil {
		return cat, 0, errDiff
	}
	return cat, lev, nil
}

func (engine InteractiveEngine) SelectWordCategory() (int, error) {
	engine.ui.SelectWordCategory()
	cat, err := engine.ui.ReadUserRune()
	if err != nil {
		return '0', err
	}
	catInt, success := domain.GetValByRune(cat)
	if !success || catInt < 0 || catInt > 2 {
		return '0', ErrBadCategory
	}
	return catInt, nil
}

func (engine InteractiveEngine) SelectLevelDifficulty() (int, error) {
	engine.ui.SelectDifficultyLevel()
	lev, err := engine.ui.ReadUserRune()
	if err != nil {
		return '0', err
	}
	levInt, success := domain.GetValByRune(lev)
	if !success || levInt < 0 || levInt > 2 {
		return '0', ErrBadDifficulty
	}
	return levInt, nil
}

func (engine InteractiveEngine) SelectAttemptCount() (int, error) {
	engine.ui.ShowAttemptSettings()
	att, err := engine.ui.ReadUserNumber()
	if err != nil {
		return 0, err
	}
	if att < 6 || att > 15 {
		return 0, errors.New("bad attempt count")
	}
	return att, nil
}

func (engine InteractiveEngine) GetRandomWordByParams() []rune {
	category, difLevel, err := engine.SelectCategoryAndDifficulty()

	switch {
	case errors.Is(err, ErrBadCategory):
		engine.ui.ShowErrorMessage("Category select automatically", err)
		category = engine.rnd.GetRandomCategory()
	case errors.Is(err, ErrBadDifficulty):
		engine.ui.ShowErrorMessage("Difficulty select automatically", err)
		difLevel = engine.rnd.GetRandomDifficulty()
	case errors.Is(err, ErrBadDifficultyAndCategory):
		engine.ui.ShowErrorMessage("Difficulty and category select automatically", err)
		category = engine.rnd.GetRandomCategory()
		difLevel = engine.rnd.GetRandomDifficulty()
	default:
	}

	engine.ui.ShowWordCategory(category)
	engine.ui.ShowDifficultyLevel(difLevel)
	hiddenWord, _ := engine.rnd.GetRandomWord(category, difLevel)
	return hiddenWord
}
