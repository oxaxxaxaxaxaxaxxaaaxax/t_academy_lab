package adapter

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"hw4-fractal-flame/internal/domain"
	"hw4-fractal-flame/internal/domain/dto"

	"github.com/stretchr/testify/assert"
)

func TestParseVariations(t *testing.T) {
	tests := []struct {
		name        string
		functions   string
		expected    domain.Variations
		expectedErr error
	}{
		{
			name:      "single function",
			functions: "heart:0.7",
			expected: domain.Variations{
				"heart": 0.7,
			},
			expectedErr: nil,
		},
		{
			name:      "multiple functions",
			functions: "heart:0.7,linear:0.3",
			expected: domain.Variations{
				"heart":  0.7,
				"linear": 0.3,
			},
			expectedErr: nil,
		},
		{
			name:        "invalid weight",
			functions:   "heart:AAA",
			expectedErr: ErrInvalidVariations,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := domain.FractalConfig{Functions: test.functions}

			variations, err := ParseVariations(cfg)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expected, variations)
		})
	}
}

func TestParseUserConfig(t *testing.T) {
	tests := []struct {
		name         string
		json         string
		expectedConf dto.UserConfig
		expectedErr  error
	}{
		{
			name: "valid config",
			json: `{
					"size": { "width": 1920, "height": 1080 },
 					"iteration_count": 300000,
					"output_path": "heart_config.png",
					"threads": 4,
					"seed": 77,
					"functions": [
						{ "name": "heart", "weight": 0.7 }
					],
					"affine_params": [
						{ "a": 0.5, "b": 1.0, "c": 1.0, "d": 1.0, "e": 1.0, "f": 1.0 }
					]
				   }`,
			expectedConf: dto.UserConfig{
				Size:           dto.Size{Width: 1920, Height: 1080},
				IterationCount: 300000,
				OutputPath:     "heart_config.png",
				Threads:        4,
				Seed:           77,
				Functions: []dto.Function{
					{Name: "heart", Weight: 0.7},
				},
				AffineParams: []dto.AffineParams{
					{A: 0.5, B: 1.0, C: 1.0, D: 1.0, E: 1.0, F: 1.0},
				},
			},
		},
		{
			name:         "invalid config",
			json:         `{ invalid }`,
			expectedErr:  ErrInvalidJSONFile,
			expectedConf: dto.UserConfig{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmp := t.TempDir()
			filePath := filepath.Join(tmp, "config.json")
			err := os.WriteFile(filePath, []byte(test.json), os.ModePerm)
			assert.NoError(t, err)

			userCfg, err := ParseUserConfig(filePath)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedConf, userCfg)
		})
	}
}

func TestValidateFractalConfig(t *testing.T) {
	tests := []struct {
		name         string
		inputConf    domain.FractalConfig
		expectedConf domain.FractalConfig
	}{
		{
			name: "fill defaults",
			inputConf: domain.NewFractalConfig(
				0, 0, 0, 0, 0,
				-1, "", "", "",
				"", "", 0,
			),
			expectedConf: domain.NewFractalConfig(
				DefaultWidth,
				DefaultHeight,
				DefaultIterationCount,
				DefaultThreads,
				DefaultSymmetryLevel,
				DefaultSeed,
				DefaultOutputPath,
				DefaultAffineParams,
				DefaultFunctions,
				"",
				DefaultIteratedAlgorithm,
				DefaultGamma,
			),
		},
		{
			name: "fill users params",
			inputConf: domain.NewFractalConfig(
				111, 222, 333, 4, 9,
				7, "x.png", "1,2,3,4,5,6", "linear:1.0",
				"cfg.json", "prob_iter", 1.7,
			),
			expectedConf: domain.NewFractalConfig(
				111, 222, 333, 4, 9,
				7, "x.png", "1,2,3,4,5,6", "linear:1.0",
				"cfg.json", "prob_iter", 1.7,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf := ValidateFractalConfig(test.inputConf)
			assert.Equal(t, test.expectedConf, conf)
		})
	}
}

func TestUpdateGameConfig(t *testing.T) {
	tests := []struct {
		name     string
		userCfg  dto.UserConfig
		cfg      domain.FractalConfig
		expected domain.FractalConfig
	}{
		{
			name: "takes values from user config",
			userCfg: dto.UserConfig{
				Size:           dto.Size{Width: 100, Height: 50},
				IterationCount: 777,
				Threads:        4,
				Seed:           3,
				OutputPath:     "out.png",
				Gamma:          1.5,
				SymmetryLevel:  6,
				Functions:      []dto.Function{{Name: "heart", Weight: 0.7}},
				AffineParams:   []dto.AffineParams{{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6}},
			},
			cfg: domain.NewFractalConfig(
				0, 0, 0, 0, 0,
				0, "", "", "",
				"", "", 0,
			),
			expected: domain.NewFractalConfig(
				100, 50, 777, 4, 6,
				3, "out.png",
				AffinesToString([]dto.AffineParams{{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6}}),
				FunctionsToString([]dto.Function{{Name: "heart", Weight: 0.7}}),
				"", "", 1.5,
			),
		},
		{
			name: "does not take values from user config",
			userCfg: dto.UserConfig{
				Size:           dto.Size{Width: 100, Height: 50},
				IterationCount: 777,
				Threads:        4,
				Seed:           3,
				OutputPath:     "out.png",
				Gamma:          1.5,
				SymmetryLevel:  6,
				Functions:      []dto.Function{{Name: "heart", Weight: 0.7}},
				AffineParams:   []dto.AffineParams{{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6}},
			},
			cfg: domain.NewFractalConfig(
				111, 222, 333, 2, 9,
				7, "x.png", "1,2,3,4,5,6", "linear:1.0",
				"", "prob_iter", 2.2,
			),
			expected: domain.NewFractalConfig(
				111, 222, 333, 2, 9,
				7, "x.png", "1,2,3,4,5,6", "linear:1.0",
				"", "prob_iter", 2.2,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := UpdateGameConfig(test.userCfg, test.cfg)
			assert.Equal(t, test.expected, cfg)
		})
	}
}

func TestFunctionsToString(t *testing.T) {
	tests := []struct {
		name     string
		dto      []dto.Function
		expected string
	}{
		{
			name: "single function",
			dto: []dto.Function{
				{Name: "heart", Weight: 0.7},
			},
			expected: "heart:0.7",
		},
		{
			name: "multiple functions",
			dto: []dto.Function{
				{Name: "heart", Weight: 0.7},
				{Name: "linear", Weight: 1.0},
			},
			expected: "heart:0.7,linear:1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			funks := FunctionsToString(test.dto)
			assert.Equal(t, test.expected, funks)
		})
	}
}

func TestAffinesToString(t *testing.T) {
	tests := []struct {
		name     string
		dto      []dto.AffineParams
		expected string
	}{
		{
			name: "single affine",
			dto: []dto.AffineParams{
				{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6},
			},
			expected: "0.1,0.2,0.3,0.4,0.5,0.6",
		},
		{
			name: "multiple affines",
			dto: []dto.AffineParams{
				{A: 0.1, B: 0.2, C: 0.3, D: 0.4, E: 0.5, F: 0.6},
				{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0},
			},
			expected: "0.1,0.2,0.3,0.4,0.5,0.6/1,0,0,0,1,0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			affines := AffinesToString(test.dto)
			assert.Equal(t, test.expected, affines)
		})
	}
}

func TestParseUserInput(t *testing.T) {
	tests := []struct {
		name           string
		cmdArgs        []string
		json           string
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "no config",
			cmdArgs:        []string{"cmd"},
			json:           "",
			expectedWidth:  DefaultWidth,
			expectedHeight: DefaultHeight,
		},
		{
			name:    "use config",
			cmdArgs: []string{"cmd", "--config", "config.json"},
			json: `{
					"size": { "width": 100, "height": 50 },
					"iteration_count": 777,
					"output_path": "out.png",
					"threads": 4,
					"seed": 3,
					"gamma": 1.5,
					"symmetry_level": 6,
					"functions": [
						{ "name": "heart", "weight": 0.7 }
					],
					"affine_params": 
						[{ "a": 0.1, "b": 0.2, "c": 0.3, "d": 0.4, "e": 0.5, "f": 0.6 }]
				   }`,
			expectedWidth:  100,
			expectedHeight: 50,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmp := t.TempDir()
			var filePath string
			var testArgs []string
			if test.json != "" {
				filePath = filepath.Join(tmp, "config.json")
				err := os.WriteFile(filePath, []byte(test.json), os.ModePerm)
				assert.NoError(t, err)

			}

			for _, arg := range test.cmdArgs {
				if arg == "config.json" {
					testArgs = append(testArgs, filePath)
					continue
				}
				testArgs = append(testArgs, arg)

			}
			CmdFlagHelper(t, testArgs, func() {

				cfg, err := ParseUserInput()
				assert.NoError(t, err)

				assert.Equal(t, test.expectedWidth, cfg.Width)
				assert.Equal(t, test.expectedHeight, cfg.Height)
			})

		})
	}
}

func TestParseCmdLine(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected domain.FractalConfig
	}{
		{
			name: "parses short flags",
			args: []string{
				"cmd",
				"-w", "100",
				"-h", "50",
				"-i", "777",
				"-t", "4",
				"-s", "6",
				"--seed", "3",
				"-o", "out.png",
				"-ap", "0.1,0.2,0.3,0.4,0.5,0.6",
				"-f", "heart:0.7",
				"--config", "cfg.json",
				"--iterated_algorithm", "prob_iter",
				"--gamma", "1.5",
			},
			expected: domain.NewFractalConfig(
				100, 50, 777, 4, 6,
				3, "out.png",
				"0.1,0.2,0.3,0.4,0.5,0.6",
				"heart:0.7",
				"cfg.json",
				"prob_iter",
				1.5,
			),
		},
		{
			name: "parses long flags",
			args: []string{
				"cmd",
				"--width", "111",
				"--height", "222",
				"--iteration-count", "333",
				"--threads", "2",
				"--symmetry-level", "5",
				"--seed", "9",
				"--output-path", "x.png",
				"--affine-params", "1,2,3,4,5,6",
				"--functions", "linear:1",
				"--iterated_algorithm", "sum_iter",
				"--gamma", "2.2",
			},
			expected: domain.NewFractalConfig(
				111, 222, 333, 2, 5,
				9, "x.png",
				"1,2,3,4,5,6",
				"linear:1",
				"",
				"sum_iter",
				2.2,
			),
		},
		{
			name: "no flags",
			args: []string{
				"cmd",
			},
			expected: domain.NewFractalConfig(
				0, 0, 0, 0, 0,
				0, "",
				"",
				"",
				"",
				"",
				0,
			),
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
