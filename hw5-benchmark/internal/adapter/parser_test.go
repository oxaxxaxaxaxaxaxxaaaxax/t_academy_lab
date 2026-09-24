package adapter

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		expectedErr error
	}{
		{
			name:        "json is valid",
			format:      "json",
			expectedErr: nil,
		},
		{
			name:        "text is valid",
			format:      "text",
			expectedErr: nil,
		},
		{
			name:        "invalid format",
			format:      "xml",
			expectedErr: ErrInvalidFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateFormat(test.format)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestValidatePath_ValidPath(t *testing.T) {
	tmp := t.TempDir()
	dirPath := filepath.Join(tmp, "some_package")
	os.Mkdir(dirPath, os.ModePerm)

	err := os.WriteFile(filepath.Join(dirPath, "file.txt"), []byte("x"), os.ModePerm)
	assert.NoError(t, err)

	_, err = ValidatePath(dirPath)

	assert.NoError(t, err)
}

func TestValidatePath_NotDirPath(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "file.txt")
	err := os.WriteFile(filePath, []byte("x"), os.ModePerm)
	assert.NoError(t, err)

	_, err = ValidatePath(filePath)

	assert.ErrorIs(t, err, ErrDirectoryNotFound)
}

func TestValidatePath_DirNotExist(t *testing.T) {
	tmp := t.TempDir()
	dirPath := filepath.Join(tmp, "dir_not_exist")

	_, err := ValidatePath(dirPath)

	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestValidateArgs(t *testing.T) {
	tmp := t.TempDir()
	absTmp, err := filepath.Abs(tmp)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		cfg         domain.InputConfig
		expectedCfg domain.InspectConfig
		expectedErr error
	}{
		{
			name: "valid args",
			cfg: domain.NewInputConfig(
				tmp,
				"Manager",
				"text",
				"out.txt",
			),
			expectedCfg: domain.NewInspectConfig(absTmp, "Manager", "text"),
			expectedErr: nil,
		},
		{
			name: "invalid format",
			cfg: domain.NewInputConfig(
				tmp,
				"Manager",
				"xml",
				"out.txt",
			),
			expectedCfg: domain.InspectConfig{},
			expectedErr: ErrInvalidFormat,
		},
		{
			name: "invalid path",
			cfg: domain.NewInputConfig(
				filepath.Join(tmp, "nope"),
				"Manager",
				"text",
				"out.txt",
			),
			expectedCfg: domain.InspectConfig{},
			expectedErr: os.ErrNotExist,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inspectCfg, err := ValidateArgs(test.cfg)

			assert.ErrorIs(t, err, test.expectedErr)
			assert.Equal(t, inspectCfg, test.expectedCfg)
		})
	}
}

func TestParseCmdLine(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected domain.InputConfig
	}{
		{
			name: "parses long flags",
			args: []string{
				"cmd",
				"-package", "./example",
				"-struct", "Manager",
				"-format", "text",
				"-output", "text_info",
			},
			expected: domain.NewInputConfig("./example", "Manager", "text", "text_info"),
		},
		{
			name: "parses short flags",
			args: []string{
				"cmd",
				"-p", "./example",
				"-s", "Manager",
				"-f", "json",
				"-output", "out.json",
			},
			expected: domain.NewInputConfig("./example", "Manager", "json", "out.json"),
		},
		{
			name: "no flags",
			args: []string{
				"cmd",
			},
			expected: domain.NewInputConfig("", "", "", ""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			CmdFlagHelper(t, test.args, func() {
				cfg := ParseCmdLine()
				assert.Equal(t, test.expected, cfg)
			})
		})
	}
}

func CmdFlagHelper(t *testing.T, args []string, fn func()) {
	t.Helper()

	oldArgs := os.Args
	oldCmd := flag.CommandLine

	flag.CommandLine = flag.NewFlagSet(args[0], flag.ContinueOnError)
	os.Args = args

	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCmd
	})

	fn()
}
