package application

import "handman/internal/domain"

//go:generate mockgen -source=user_interface.go -destination=../mocks/mock_user_interface.go -package=mocks

type UI interface {
	ShowMatchResultWord(word string, stat domain.Status)
	ShowWord(word string)
	ShowBadInput()
	ShowLose()
	ShowHangmanPicture(picture string)
	ShowWin()
	ShowWordCategory(category int)
	SelectWordCategory()
	SelectDifficultyLevel()
	ShowDifficultyLevel(diffLevel int)
	ShowAttemptSettings()
	ShowError(e error)
	ShowCountMistakes(countAttempts, CountMistakes int)
	ShowErrorMessage(message string, err error)
	ShowDescription(desc string)
	ReadUserAnswer() string
	ReadHiddenWord() string
	GetLengthInput() int
	ReadUserRune() (rune, error)
	ReadUserNumber() (int, error)
}
