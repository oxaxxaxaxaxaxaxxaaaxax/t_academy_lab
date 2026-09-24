package application

import (
	"errors"
	"hw4-fractal-flame/internal/domain"
	"log/slog"
	"strconv"
	"strings"
)

var (
	ErrInvalidConfig         = errors.New("invalid config")
	ErrWrongAffineParameters = errors.New("wrong affine parameters")
)

type AffineParams struct {
	A float64
	B float64
	C float64
	D float64
	E float64
	F float64
}

type CoefficientArray struct {
	AffineCoefficients []AffineParams
	N                  int
}

func ParseAffineParams(cfg domain.FractalConfig) (CoefficientArray, error) {
	var coArray CoefficientArray
	var err error

	setsParam := strings.Split(cfg.AffineParams, "/")
	coArray.N = len(setsParam)
	coArray.AffineCoefficients = make([]AffineParams, coArray.N)
	coArr := coArray.AffineCoefficients

	for eqIdx, params := range setsParam {
		coefficients := strings.Split(params, ",")
		if len(coefficients) != 6 {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}
		coArr[eqIdx].A, err = strconv.ParseFloat(coefficients[0], 64)
		if err != nil {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}

		coArr[eqIdx].B, err = strconv.ParseFloat(coefficients[1], 64)
		if err != nil {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}

		coArr[eqIdx].C, err = strconv.ParseFloat(coefficients[2], 64)
		if err != nil {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}

		coArr[eqIdx].D, err = strconv.ParseFloat(coefficients[3], 64)
		if err != nil {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}

		coArr[eqIdx].E, err = strconv.ParseFloat(coefficients[4], 64)
		if err != nil {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}

		coArr[eqIdx].F, err = strconv.ParseFloat(coefficients[5], 64)
		if err != nil {
			return CoefficientArray{}, errors.Join(ErrInvalidConfig, ErrWrongAffineParameters)
		}
	}

	slog.Debug("affine parameters", "params", coArr)
	coArray.AffineCoefficients = coArr

	return coArray, nil
}
