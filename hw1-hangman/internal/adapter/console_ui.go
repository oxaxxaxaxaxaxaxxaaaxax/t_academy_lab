package adapter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"handman/internal/domain"
)

type ConsoleUI struct {
	in *bufio.Reader
}

func NewConsoleUI(r io.Reader) ConsoleUI {
	return ConsoleUI{in: bufio.NewReader(r)}
}

func (ConsoleUI) ShowMatchResultWord(word string, stat domain.Status) {
	if stat == domain.Success {
		fmt.Print(word, ";", "POS")
	} else {
		fmt.Print(word, ";", "NEG")
	}

	ch := make(chan rune)
	close(ch)
}

func (ConsoleUI) ShowWord(word string) {
	fmt.Println(word)
}

func (ConsoleUI) ShowBadInput() {
	fmt.Println("Bad input, try again")
}

func (ConsoleUI) ShowLose() {
	fmt.Println("You lose")
}

func (ConsoleUI) ShowHangmanPicture(picture string) {
	fmt.Println(picture)
}

func (ConsoleUI) ShowWin() {
	fmt.Println("You win!")
}

func (ConsoleUI) SelectWordCategory() {
	fmt.Println("Select category")
	for i := range domain.CategoryCount() {
		fmt.Println(domain.GetCategoryName(i))
	}
}

func (ConsoleUI) ShowWordCategory(category int) {
	fmt.Println("Category name: ", category)
}

func (ConsoleUI) SelectDifficultyLevel() {
	fmt.Println("Select category")
	for i := range domain.DifficultyCount() {
		fmt.Println(domain.GetLevelDifficulty(i))
	}
}

func (ConsoleUI) ShowDifficultyLevel(difficultyLevel int) {
	fmt.Println("Difficulty level: ", difficultyLevel)
}

func (ConsoleUI) ShowAttemptSettings() {
	fmt.Println("Select the number of attempts (maximum 15): Current number is 6")
}

func (ConsoleUI) ShowError(e error) {
	fmt.Println(e)
}

func (ConsoleUI) ShowCountMistakes(countAttempts, CountMistakes int) {
	fmt.Printf("You have: %d attempts", countAttempts-CountMistakes)
}

func (ConsoleUI) ShowErrorMessage(message string, err error) {
	fmt.Println(message, ", cause:", err.Error())
}

func (ConsoleUI) ShowDescription(desc string) {
	fmt.Println("help:", desc)
}

func (ConsoleUI) ReadUserAnswer() string {
	return os.Args[2]
}

func (ConsoleUI) ReadHiddenWord() string {
	return os.Args[1]
}

func (ConsoleUI) GetLengthInput() int {
	return len(os.Args)
}

func (ui ConsoleUI) ReadUserRune() (rune, error) {
	for {
		str, err := ui.in.ReadString('\n')
		if err != nil {
			return 0, err
		}
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		rns := []rune(str)
		if len(rns) == 1 {
			return rns[0], nil
		}
		fmt.Println("Invalid input, try again")
		continue
	}
}

func (ui ConsoleUI) ReadUserNumber() (int, error) {
	n := 0
	digits := 0

	for {
		rn, _, err := ui.in.ReadRune()
		if err != nil {
			return 0, err
		}
		if rn == '\n' && digits > 0 {
			return n, nil
		}

		if unicode.IsSpace(rn) {
			continue
		}
		n = n*10 + int(rn-'0')
		digits++
	}
}
