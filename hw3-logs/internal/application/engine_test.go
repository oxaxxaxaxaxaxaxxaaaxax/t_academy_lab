package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

func TestSaveStats(t *testing.T) {
	type testCase struct {
		conf domain.Config
		dto  models.DTO
	}
	testCases := []testCase{
		{
			conf: domain.Config{Output: "../../scripts/data/output/data.json", Format: "json"},
			dto:  models.DTO{},
		},
		{
			conf: domain.Config{Output: "../../scripts/data/output/data.md", Format: "markdown"},
			dto:  models.DTO{},
		},
		{
			conf: domain.Config{Output: "../../scripts/data/output/data.ad", Format: "adoc"},
			dto:  models.DTO{},
		},
	}

	for _, test := range testCases {
		err := SaveStats(test.dto, test.conf)
		assert.Nil(t, err)
	}
}
