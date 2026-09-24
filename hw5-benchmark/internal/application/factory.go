package application

import (
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

func NewFormatter(format string) domain.Formatter {
	switch strings.ToLower(format) {
	case "json":
		return NewJsonFormatter()
	case "text":
		return NewTextFormatter()
	default:
		return NewTextFormatter()
	}
}
