package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

func TestNewScanner(t *testing.T) {
	type testCase struct {
		sourceType       SourceType
		expectedErr      error
		expectedScanType domain.LogScanner
	}

	testCases := []testCase{
		{
			sourceType:       Local,
			expectedErr:      nil,
			expectedScanType: LocalLogScanner{},
		},
		{
			sourceType:       Remote,
			expectedErr:      nil,
			expectedScanType: RemoteLogScanner{},
		},
	}

	for _, test := range testCases {
		scan, err := NewScanner(test.sourceType)
		assert.Equal(t, test.expectedErr, err)
		assert.IsType(t, test.expectedScanType, scan)
	}
}
