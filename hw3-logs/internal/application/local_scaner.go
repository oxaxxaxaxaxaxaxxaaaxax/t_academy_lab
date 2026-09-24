package application

import (
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

const TimeLayout string = "2/Jan/2006:15:04:05 -0700"

type LocalLogScanner struct {
}

func NewLocalLogScanner() LocalLogScanner {
	return LocalLogScanner{}
}

func (l LocalLogScanner) StartScanLogs(ctx domain.LogCtx) (models.DTO, error) {
	slog.Info("start scan logs in local files")
	paths := ctx.FilePath
	logStats := NewLogStats(paths)
	var err error
	var file *os.File

	for _, path := range paths {
		file, err = os.Open(path)
		if err != nil {
			return models.DTO{}, err
		}
		defer file.Close()

		logStats, err = logStats.GetLogStats(file)
		if err != nil {
			return models.DTO{}, ErrLogFormatMismatch
		}
	}

	logStats, err = logStats.CalculateRemainingStats()
	if err != nil {
		return models.DTO{}, err
	}

	dto := logStats.ConvertToDto()

	return dto, nil
}
