package application

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {
	type TestCase struct {
		filePath       string
		expectedErr    error
		expectedFormat SourceType
	}

	tests := []TestCase{
		{
			//I don't know how to test this case without relative path, I use files from blackbox tests
			filePath:       "../../scripts/data/input/logs/*",
			expectedErr:    nil,
			expectedFormat: Local,
		},
		{
			filePath:       "../../scripts/data/input/logs/part1.txt",
			expectedErr:    nil,
			expectedFormat: Local,
		},
		{
			filePath:       "https://raw.githubusercontent.com/elastic/examples/master/Common%20Data%20Formats/nginx_logs/nginx_logs",
			expectedErr:    nil,
			expectedFormat: Remote,
		},
	}

	for _, test := range tests {
		logSource, err := ParsePath(test.filePath)
		assert.Equal(t, test.expectedErr, err)
		assert.Equal(t, test.expectedFormat, logSource.SourceType)
	}
}

func TestGetLogStats(t *testing.T) {
	type TestCase struct {
		filePath    string
		expectedErr error
	}

	tests := []TestCase{
		{
			filePath:    "../../scripts/data/input/logs/part1.txt",
			expectedErr: nil,
		},
		{
			filePath:    "../../scripts/data/input/file2.txt",
			expectedErr: ErrLogFormatMismatch,
		},
	}
	for _, test := range tests {
		file, err := os.Open(test.filePath)
		assert.Nil(t, err)
		defer file.Close()
		l := NewLogStats([]string{test.filePath})
		_, err = l.GetLogStats(file)
		assert.Equal(t, test.expectedErr, err)
	}
}
