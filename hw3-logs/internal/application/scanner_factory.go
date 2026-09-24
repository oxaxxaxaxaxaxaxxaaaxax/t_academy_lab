package application

import (
	"github.com/pkg/errors"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

var ErrInvalidSourceType = errors.New("invalid source type")

func NewScanner(sourceType SourceType) (domain.LogScanner, error) {
	switch sourceType {
	case Local:
		return NewLocalLogScanner(), nil
	case Remote:
		return NewRemoteLogScanner(), nil
	}
	return LocalLogScanner{}, ErrInvalidSourceType
}
