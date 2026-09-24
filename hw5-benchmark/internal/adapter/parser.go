package adapter

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

var (
	ErrInvalidFormat     = errors.New("invalid format")
	ErrDirectoryNotFound = errors.New("directory not found")
)

func ParseCmdLine() domain.InputConfig {
	pkgName := flag.String("package", "", "")
	flag.StringVar(pkgName, "p", "", "")

	structName := flag.String("struct", "", "")
	flag.StringVar(structName, "s", "", "")

	format := flag.String("format", "", "")
	flag.StringVar(format, "f", "", "")

	output := flag.String("output", "", "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
		Usage:
		  inspector [options]
		
		Options:
		  -package, -p string
				go package, in which struct is located
		
		  -struct, -s string
				struct name to analize
		
		  -format, -f string
				file format to output result
		
		Examples:
		  ./main -package ./example -struct Manager -format text -output text_info
		`)
	}

	flag.Parse()

	return domain.NewInputConfig(*pkgName, *structName, *format, *output)

}

func ValidateArgs(cfg domain.InputConfig) (domain.InspectConfig, error) {
	err := ValidateFormat(cfg.Format)
	if err != nil {
		return domain.InspectConfig{}, err
	}
	absPkgPath, err := ValidatePath(cfg.PkgName)
	if err != nil {
		return domain.InspectConfig{}, err
	}
	return domain.NewInspectConfig(absPkgPath, cfg.StructName, cfg.Format), nil
}

func ValidateFormat(format string) error {
	if format != "json" && format != "text" {
		return ErrInvalidFormat
	}
	return nil
}

func ValidatePath(packagePath string) (string, error) {
	absPath, err := filepath.Abs(packagePath)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", ErrDirectoryNotFound
	}
	return absPath, nil
}
