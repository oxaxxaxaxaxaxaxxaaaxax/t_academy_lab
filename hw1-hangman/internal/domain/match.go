package domain

import (
	"unicode"
)

type Status int

const (
	Success Status = iota
	Fail
)

func MatchWords(hiddenWord, userAnswer string) (string, Status) {
	runes1 := []rune(hiddenWord)
	runes2 := []rune(userAnswer)
	lenStr := len(runes1)
	rs := []rune(userAnswer)

	for i := 0; i < lenStr; i++ {
		_, stat := MatchLetter(runes2[i], userAnswer, hiddenWord)
		if stat == Fail {
			rs[i] = '*'
		}
	}

	if IsEqualWords(hiddenWord, string(rs)) || IsContainsAllLetters(string(rs)) {
		return hiddenWord, Success
	}
	return string(rs), Fail
}

func IsContainsAllLetters(str string) bool {
	for _, let := range str {
		if rune(let) == '*' {
			return false
		}
	}
	return true
}

func IsEqualWords(str1, str2 string) bool {
	equal := true
	runes1 := []rune(str1)
	runes2 := []rune(str2)
	result := []rune{}
	lenStr := len(runes1)

	for i := 0; i < lenStr; i++ {
		if unicode.ToLower(runes1[i]) == unicode.ToLower(runes2[i]) {
			result = append(result, runes1[i])
		} else {
			result = append(result, '*')
			equal = false
		}
	}
	return equal
}

func MatchLetter(letter rune, current, original string) (string, Status) {
	found := false
	runes := []rune(original)
	result := []rune(current)
	lenStr := len(runes)

	for i := 0; i < lenStr; i++ {
		if unicode.ToLower(runes[i]) == unicode.ToLower(letter) {
			result[i] = unicode.ToLower(letter)
			found = true
		}
	}
	if found {
		return string(result), Success
	}
	return string(result), Fail
}

func CountDiffSymbols(a, b []rune) int {
	counter := 0

	for i := range len(a) {
		if a[i] != b[i] {
			counter++
		}
	}
	return counter
}

func CurrentPictureIndex(mistakes, countPictures, countAttempt int) int {
	return mistakes * (countPictures - 1) / countAttempt
}
