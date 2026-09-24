package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFilePath(t *testing.T) {
	type TestCase struct {
		filePath string
		expected error
	}

	testCases := []TestCase{
		{
			//I don't know how to test this case without relative path, I use files from blackbox tests
			filePath: "../../scripts/data/input/logs/part1.txt",
			expected: nil,
		},
		{
			filePath: "oaoaoaoaoalog",
			expected: ErrInvalidPathToFile,
		},
		{
			filePath: "oaoaoaoao.txt",
			expected: ErrFileNotFound,
		},
		{
			filePath: "https://oaoaoaoa.txt",
			expected: ErrInvalidPathToFile,
		},
	}

	for _, testCase := range testCases {
		err := CheckFilePathValidate(testCase.filePath)
		assert.Equal(t, testCase.expected, err)
	}
}

func TestCheckFormatValidate(t *testing.T) {
	type TestCase struct {
		fileFormat string
		expected   error
	}

	testCases := []TestCase{
		{
			fileFormat: "json",
			expected:   nil,
		},
		{
			fileFormat: "markdown",
			expected:   nil,
		},
		{
			fileFormat: "adoc",
			expected:   nil,
		},
		{
			fileFormat: "tuutu",
			expected:   ErrInvalidFileFormat,
		},
	}

	for _, testCase := range testCases {
		err := CheckFormatValidate(testCase.fileFormat)
		assert.Equal(t, testCase.expected, err)
	}
}

func TestCheckFileExtension(t *testing.T) {
	type TestCase struct {
		fileExtension string
		outputPath    string
		expected      error
	}

	testCases := []TestCase{
		{
			fileExtension: "json",
			outputPath:    "logs.json",
			expected:      nil,
		},
		{
			fileExtension: "markdown",
			outputPath:    "logs.md",
			expected:      nil,
		},
		{
			fileExtension: "json",
			outputPath:    "logs.ad",
			expected:      ErrInvalidFileExtension,
		},
		{
			fileExtension: "adoc",
			outputPath:    "logs.ad",
			expected:      nil,
		},
	}

	for _, testCase := range testCases {
		err := CheckExtension(testCase.outputPath, testCase.fileExtension)
		assert.Equal(t, testCase.expected, err)
	}
}

func TestCheckFileExist(t *testing.T) {
	type TestCase struct {
		filePath string
		expected error
	}

	testCases := []TestCase{
		{
			filePath: "../../scripts/data/input/logs/part1.txt",
			expected: ErrOutputFileAlreadyExists,
		},
		{
			filePath: "../not_exist_file.md",
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		err := CheckFileExist(testCase.filePath)
		assert.Equal(t, testCase.expected, err)
	}
}

func TestCheckOrderDate(t *testing.T) {
	type TestCase struct {
		fileFrom string
		fileTo   string
		expected error
	}

	testCases := []TestCase{
		{
			fileFrom: "2025-01-02",
			fileTo:   "2025-01-01",
			expected: ErrInvalidDateOrder,
		},
		{
			fileFrom: "2025-01-01",
			fileTo:   "-",
			expected: nil,
		},
		{
			fileFrom: "-",
			fileTo:   "2025-01-01",
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		err := CheckOrderDate(testCase.fileFrom, testCase.fileTo)
		assert.Equal(t, testCase.expected, err)
	}
}
