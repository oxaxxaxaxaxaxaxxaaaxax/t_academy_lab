package domain

import (
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

type LogScanner interface {
	StartScanLogs(ctx LogCtx) (models.DTO, error)
}

type PathChecker interface {
	Check(filePath []string) error
}
