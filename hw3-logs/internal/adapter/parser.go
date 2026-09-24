package adapter

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
)

func ParseCmdLine() domain.Config {
	path := flag.String("p", "", "path to nginx files")
	format := flag.String("f", "json", "format to output results")
	from := flag.String("from", "-", "from date")
	to := flag.String("to", "-", "to date")
	output := flag.String("o", "console", "file to save results")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
		Help message for usage:

		-p
			Write path to nginx files for logging
		-f
			Format to output file with results"
		-o
			File to save log-results
		-from
			Date from (format ISO8601)
		-to
			Date to (format ISO8601)`)
	}
	flag.Parse()
	slog.Info("Parsing flags",
		"path", *path,
		"format", *format,
		"from", *from,
		"to", *to,
		"output", *output,
	)
	return domain.Config{PathToFile: *path, Format: *format,
		Output: *output, From: *from, To: *to}
}
