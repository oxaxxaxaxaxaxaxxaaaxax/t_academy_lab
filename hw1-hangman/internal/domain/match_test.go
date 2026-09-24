package domain

import (
	"testing"
)

func TestMatchWords(t *testing.T) {
	type ts struct {
		word1        string
		word2        string
		expectedWord string
		expectedStat Status
	}

	cases := []ts{
		{
			word1: "волокно\n", word2: "толокно\n", expectedWord: "*олокно\n", expectedStat: Fail,
		},
		{
			word1: "кулон\n", word2: "клоун\n", expectedWord: "кулон\n", expectedStat: Success,
		},
	}

	for _, val := range cases {
		res, stat := MatchWords(val.word1, val.word2)
		if res != val.expectedWord || stat != val.expectedStat {
			t.Errorf("Failed MatchWords: got %v,%v want %v,%v", res, stat, val.expectedWord, val.expectedStat)
		}
	}
}

func TestGetDiffSymbols(t *testing.T) {
	type ts struct {
		word1         []rune
		word2         []rune
		expectedCount int
	}

	cases := []ts{
		{
			word1: []rune("space"), word2: []rune("spice"), expectedCount: 1,
		},
		{
			word1: []rune(""), word2: []rune(""), expectedCount: 0,
		},
	}

	for _, val := range cases {
		if counter := CountDiffSymbols(val.word1, val.word2); counter != val.expectedCount {
			t.Errorf("GetDiffSymbols failed: got %v want %v", counter, val.expectedCount)
		}
	}
}
