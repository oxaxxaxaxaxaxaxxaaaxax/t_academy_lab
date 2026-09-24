package application

import (
	"log/slog"
	"net/http"
	"time"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

const (
	timeTimeout int = 5
)

type RemoteLogScanner struct{}

func NewRemoteLogScanner() RemoteLogScanner {
	return RemoteLogScanner{}
}

func (r RemoteLogScanner) StartScanLogs(ctx domain.LogCtx) (models.DTO, error) {
	slog.Info("start scan logs in remote file")
	client := http.Client{Timeout: time.Duration(timeTimeout) * time.Second}
	httpFilePath := ctx.FilePath[0]

	httpReq, err := http.NewRequest(http.MethodGet, httpFilePath, nil)
	if err != nil {
		return models.DTO{}, err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return models.DTO{}, err
	}

	defer resp.Body.Close()

	logStats := NewLogStats([]string{httpFilePath})
	logStats, err = logStats.GetLogStats(resp.Body)
	if err != nil {
		return models.DTO{}, err
	}
	logStats, err = logStats.CalculateRemainingStats()
	if err != nil {
		return models.DTO{}, err
	}

	dto := logStats.ConvertToDto()

	return dto, nil
}
