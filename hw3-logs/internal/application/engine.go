package application

import (
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

func StartApplication(config domain.Config) error {
	log, err := ParsePath(config.PathToFile)
	if err != nil {
		slog.Warn("Parsing error:", "err", err)
		return err
	}
	scanner, err := NewScanner(log.SourceType)
	if err != nil {
		slog.Warn("Scanner factory finished with error:", "err", err)
		return err
	}

	ctx := domain.LogCtx{FilePath: log.Path}
	dto, err := scanner.StartScanLogs(ctx)
	if err != nil {
		slog.Warn("Scan logs finished with error:", "err", err)
		return err
	}

	err = SaveStats(dto, config)
	if err != nil {
		slog.Warn("Saving stats finished with error:", "err", err)
		return err
	}
	return nil
}

func SaveStats(dto models.DTO, config domain.Config) error {
	slog.Info("save stats in", "out", config.Output)
	file, err := os.Create(config.Output)
	if err != nil {
		return err
	}
	switch config.Format {
	case "json":
		err = adapter.WriteJson(file, dto)
		if err != nil {
			return err
		}
	case "markdown":
		adapter.WriteMarkdown(file, config, dto)
	case "adoc":
		adapter.WriteAdoc(file, config, dto)
	}
	return nil
}
