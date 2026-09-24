package adapter

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestConsoleUIReadUserRune(t *testing.T) {
	type ts struct {
		in           string
		expectedRune rune
	}

	cases := []ts{
		{
			in: "c\n", expectedRune: 'c',
		},
		{
			in: "    \nc\n", expectedRune: 'c',
		},
		{
			in: "kpmxoxiwoxnixsx\nc\n", expectedRune: 'c',
		},
	}

	for _, val := range cases {
		ui := NewConsoleUI(strings.NewReader(val.in))
		if rn, err := ui.ReadUserRune(); err != nil || rn != val.expectedRune {
			t.Fatalf("ReadUserRune() got %q, %v want %q,%v", rn, err, val.expectedRune, nil)
		}
	}
}

func TestConsoleUIReadUserNumber(t *testing.T) {
	type ts struct {
		in             string
		expectedNumber int
		expectedErr    error
	}

	cases := []ts{
		{
			in: "7\n", expectedNumber: 7, expectedErr: nil,
		},
		{
			in: "  767\n", expectedNumber: 767, expectedErr: nil,
		},
		{
			in: "\n8", expectedNumber: 0, expectedErr: io.EOF,
		},
	}

	for _, val := range cases {
		ui := NewConsoleUI(strings.NewReader(val.in))
		if par, err := ui.ReadUserNumber(); errors.Is(err, val.expectedErr) || par != val.expectedNumber {
			t.Fatalf("ReadUserRune() got %q, %v want %q,%v", par, err, val.expectedNumber, nil)
		}
	}
}
