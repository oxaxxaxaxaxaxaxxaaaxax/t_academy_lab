package domain

import "gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks.go -package=mocks

type Formatter interface {
	Format(info dto.StructInfo) (string, error)
}

type Random interface {
	RandBoll() bool
	RandInt64() int64
	RandIntPositive() int
	RandUint64() uint64
	RandFloat64() float64
	RandString() string
}
