package application

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

const (
	WXPermission      uint32 = 0300
	JsonExtension     string = ".json"
	MarkdownExtension string = ".md"
	AdocExtension     string = ".ad"
)

var (
	ErrInvalidPathToFile       = errors.New("invalid path to file")
	ErrFileNotFound            = errors.New("file not found")
	ErrInvalidFileFormat       = errors.New("invalid file format")
	ErrInvalidFileExtension    = errors.New("invalid file extension")
	ErrOutputFileAlreadyExists = errors.New("output file already exists")
	ErrDirectoryNotInWriteMode = errors.New("directory not in write mode")
	ErrInvalidDateFormat       = errors.New("invalid date format")
	ErrInvalidDateOrder        = errors.New("invalid date order")
)

func CheckConfigValidate(config domain.Config) error {
	slog.Info("Start checking configuration")

	err := CheckFilePathValidate(config.PathToFile)
	if err != nil {
		slog.Warn("Configuration file path is invalid")
		return err
	}

	err = CheckFormatValidate(config.Format)
	if err != nil {
		slog.Warn("Configuration file format is invalid")
		return err
	}

	err = CheckOutputValidate(config.Output, config.Format)
	if err != nil {
		slog.Warn("Configuration file output is invalid")
		return err
	}

	err = CheckOrderDate(config.From, config.To)
	if err != nil {
		slog.Warn("Configuration from-to parameters is invalid")
		return err
	}
	return nil
}

func CheckFilePathValidate(filePath string) error {
	slog.Debug("CheckFilePathValidate", "path", filePath)

	logSource, err := ParsePath(filePath)
	if err != nil {
		return ErrInvalidPathToFile
	}

	switch logSource.SourceType {
	case Remote:
		slog.Debug("remote log source")
		checker := NewRemoteChecker()
		return checker.Check(logSource.Path)
	case Local:
		slog.Debug("local log source")
		checker := NewLocalChecker()
		return checker.Check(logSource.Path)
	}
	return nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func CheckFormatValidate(format string) error {
	if format != "json" && format != "markdown" && format != "adoc" {
		return ErrInvalidFileFormat
	}
	return nil
}

func CheckOutputValidate(outputPath, format string) error {
	err := CheckExtension(outputPath, format)
	if err != nil {
		return err
	}
	err = CheckFileExist(outputPath)
	if err != nil {
		return err
	}
	err = CheckDirPermissions(outputPath)
	if err != nil {
		return err
	}
	return nil
}

func CheckExtension(outputPath, format string) error {
	extensions := CreateExtensions()
	if !strings.HasSuffix(outputPath, extensions[format]) {
		return ErrInvalidFileExtension
	}
	//if we are here we have already pass format checking then return nil is safely
	return nil
}

func CreateExtensions() map[string]string {
	extensions := make(map[string]string)

	extensions["json"] = JsonExtension
	extensions["markdown"] = MarkdownExtension
	extensions["adoc"] = AdocExtension

	return extensions
}

func CheckFileExist(filePath string) error {
	exist, _ := Exists(filePath)
	if exist {
		return ErrOutputFileAlreadyExists
	}
	return nil
}

func CheckDirPermissions(outputPath string) error {
	dirPath := filepath.Dir(outputPath)
	info, err := os.Stat(dirPath)
	if err != nil {
		return ErrInvalidFileFormat
	}
	mode := info.Mode()
	if mode&os.FileMode(WXPermission) != os.FileMode(WXPermission) {
		return ErrDirectoryNotInWriteMode
	}
	return nil
}

func CheckOrderDate(fromDate, toDate string) error {
	slog.Info("Start checking date parameters")
	slog.Debug("Start checking date order")
	var (
		from time.Time
		to   time.Time
		err  error
	)

	if fromDate != "-" {
		from, err = time.Parse(time.DateOnly, fromDate)
		if err != nil {
			return ErrInvalidDateFormat
		}
	}
	if toDate != "-" {
		to, err = time.Parse(time.DateOnly, toDate)
		if err != nil {
			return ErrInvalidDateFormat
		}
	}

	if fromDate != "-" && toDate != "-" && from.After(to) {
		return ErrInvalidDateOrder
	}
	return nil
}
