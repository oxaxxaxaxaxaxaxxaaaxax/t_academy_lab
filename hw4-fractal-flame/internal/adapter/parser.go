package adapter

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"hw4-fractal-flame/internal/domain/dto"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"hw4-fractal-flame/internal/domain"
)

const (
	DefaultWidth             = 1920
	DefaultHeight            = 1080
	DefaultIterationCount    = 2500
	DefaultThreads           = 1
	DefaultSeed              = int64(5)
	DefaultOutputPath        = "result.png"
	DefaultIteratedAlgorithm = "sum_iter"
	DefaultGamma             = 2.2
	DefaultSymmetryLevel     = 1
	DefaultAffineParams      = "0.5,1.0,1.0,1.0,1.0,1.0"
	DefaultFunctions         = "heart:1.0"

	ZeroWidth             = 0
	ZeroHeight            = 0
	ZeroIterationCount    = 0
	ZeroThreads           = 0
	ZeroSeed              = 0
	ZeroGamma             = 0
	ZeroSymmetryLevel     = 0
	ZeroConfigPath        = ""
	ZeroOutputPath        = ""
	ZeroAffineParams      = ""
	ZeroFunctions         = ""
	ZeroIteratedAlgorithm = ""
)

var (
	ErrInvalidJSONFile   = errors.New("invalid js on file")
	ErrCannotReadFile    = errors.New("cannot read file")
	ErrInvalidVariations = errors.New("invalid variations")
)

func ParseUserInput() (domain.FractalConfig, error) {
	cfg := ParseCmdLine()
	if cfg.ConfigPath != ZeroConfigPath {
		userCfg, err := ParseUserConfig(cfg.ConfigPath)
		if err != nil {
			return domain.FractalConfig{}, err
		}
		cfg = UpdateGameConfig(userCfg, cfg)
	}
	cfg = ValidateFractalConfig(cfg)
	return cfg, nil
}

func ParseCmdLine() domain.FractalConfig {
	width := flag.Int("width", ZeroWidth, "result image width")
	flag.IntVar(width, "w", ZeroWidth, "result image width (short)")

	height := flag.Int("height", ZeroHeight, "result image height")
	flag.IntVar(height, "h", ZeroHeight, "result image height (short)")

	seed := flag.Int64("seed", ZeroSeed, "seed for random generator")

	iterCount := flag.Int("iteration-count", ZeroIterationCount, "number of generation iterations")
	flag.IntVar(iterCount, "i", ZeroIterationCount, "number of generation iterations (short)")

	outputPath := flag.String("output-path", "", "path to output PNG file")
	flag.StringVar(outputPath, "o", "", "path to output PNG file (short)")

	threads := flag.Int("threads", ZeroThreads, "number of threads")
	flag.IntVar(threads, "t", ZeroThreads, "number of threads (short)")

	affineParams := flag.String("affine-params", "",
		"affine transform config: a1,b1,c1,d1,e1,f1/a2,b2,c2,d2,e2,f2/...")
	flag.StringVar(affineParams, "ap", "",
		"affine transform config (short)")

	functions := flag.String("functions", "",
		"transformation functions config: name1:weight1,name2:weight2")
	flag.StringVar(functions, "f", "",
		"transformation functions config (short)")

	configPath := flag.String("config", "", "path to optional config file")

	iter := flag.String("iterated_algorithm", "",
		"iterated function algorithm: alg_name")

	gamma := flag.Float64("gamma", ZeroGamma, "coefficient for gamma correction")

	symmetry := flag.Int("symmetry-level", ZeroSymmetryLevel, "coefficient for symmetry")
	flag.IntVar(symmetry, "s", ZeroSymmetryLevel, "coefficient for symmetry (short)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
		Usage:
		 fractal [options]
		
		Options:
		 -w, --width int
		       Output image width (default 1920)
		 -h, --height int
		       Output image height (default 1080)
		     --seed int64
		       Random generator seed (default 5)
		 -i, --iteration-count int
		       Number of generation iterations (default 2500)
		 -o, --output-path string
		       Path to the output PNG (default "result.png")
		 -t, --threads int
		       Number of threads (default 1)
		 -ap, --affine-params string
		       Affine transforms:
		       a1,b1,c1,d1,e1,f1/a2,b2,c2,d2,e2,f2/...
		 -f, --functions string
		       config transformation methods:
		       name1:weight1,name2:weight2
		     --config string
		       Path to input config file
			 --iterated_algorithm string
			   kind of iterated function algorithm for chaos game
			 --gamma
			   coefficient for gamma correction
		 -s, --symmetry-level
			   coefficient for symmetry 
		
		Examples:
		 ./fractalflame -w 1920 -h 1080 -i 800000 -o out.png -ap "0.6,-0.2,0,0.2,0.6,0/0.6,0.2,0,-0.2,0.6,0/0.4,0,0,0,0.4,-0.4" -f "heart:1.0,linear:0.3,diamond:0.5,fisheye:0.4,swirl:0.2,cosine:0.2"

 		 ./fractalflame -w 1920 -h 1080 -i 50000 -threads=4 -o heart.png --seed 3 -ap "0.5,1.0,1.0,1.0,1.0,1.0" -f "heart:0.7"`)
	}

	flag.Parse()

	slog.Debug("cfg.Functions", "slice ", *functions)

	cfg := domain.NewFractalConfig(*width, *height, *iterCount, *threads, *symmetry, *seed, *outputPath, *affineParams,
		*functions, *configPath, *iter, *gamma)
	slog.Debug("cfg", "config", cfg)
	return cfg
}

func ParseVariations(cfg domain.FractalConfig) (domain.Variations, error) {
	variations := domain.Variations{}
	var funks []string

	if !strings.Contains(cfg.Functions, ",") {
		funks = append(funks, cfg.Functions)
		slog.Debug("funks with len 1", "funk", funks)
	} else {
		funks = strings.Split(cfg.Functions, ",")
		slog.Debug("funks with len != 1", "funks ", funks)
	}

	for _, function := range funks {
		funcWeightPair := strings.Split(function, ":")
		slog.Debug("funcWeightPair", "funcWeightPair", funcWeightPair)
		slog.Info("slice function:", "func", funcWeightPair)
		weight, err := strconv.ParseFloat(funcWeightPair[1], 64)
		if err != nil {
			return nil, ErrInvalidVariations
		}
		variations[funcWeightPair[0]] = weight
	}

	return variations, nil
}

func ParseUserConfig(filepath string) (dto.UserConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return dto.UserConfig{}, ErrCannotReadFile
	}

	var uCfg dto.UserConfig
	if err = json.Unmarshal(data, &uCfg); err != nil {
		return dto.UserConfig{}, ErrInvalidJSONFile
	}

	return uCfg, nil
}

func ValidateFractalConfig(cfg domain.FractalConfig) domain.FractalConfig {
	if cfg.Width <= ZeroWidth {
		cfg.Width = DefaultWidth
	}
	if cfg.Height <= ZeroHeight {
		cfg.Height = DefaultHeight
	}
	if cfg.IterCount <= ZeroIterationCount {
		cfg.IterCount = DefaultIterationCount
	}
	if cfg.Threads <= ZeroThreads {
		cfg.Threads = DefaultThreads
	}
	if cfg.Seed < ZeroSeed {
		cfg.Seed = DefaultSeed
	}
	if cfg.GammaCoefficient <= ZeroGamma {
		cfg.GammaCoefficient = DefaultGamma
	}
	if cfg.SymmetryLevel <= ZeroSymmetryLevel {
		cfg.SymmetryLevel = DefaultSymmetryLevel
	}
	if cfg.OutputPath == ZeroOutputPath {
		cfg.OutputPath = DefaultOutputPath
	}
	if cfg.IterAlgorithmName == ZeroIteratedAlgorithm {
		cfg.IterAlgorithmName = DefaultIteratedAlgorithm
	}
	if cfg.AffineParams == ZeroAffineParams {
		cfg.AffineParams = DefaultAffineParams
	}
	if cfg.Functions == ZeroFunctions {
		cfg.Functions = DefaultFunctions
	}
	return cfg
}

func UpdateGameConfig(userCfg dto.UserConfig, cfg domain.FractalConfig) domain.FractalConfig {
	if cfg.Width == ZeroWidth {
		cfg.Width = userCfg.Size.Width
	}
	if cfg.Height == ZeroHeight {
		cfg.Height = userCfg.Size.Height
	}
	if cfg.IterCount == ZeroIterationCount {
		cfg.IterCount = userCfg.IterationCount
	}
	if cfg.Threads == ZeroThreads {
		cfg.Threads = userCfg.Threads
	}
	if cfg.Seed == ZeroSeed {
		cfg.Seed = int64(userCfg.Seed)
	}
	if cfg.OutputPath == ZeroOutputPath {
		cfg.OutputPath = userCfg.OutputPath
	}
	if cfg.GammaCoefficient == ZeroGamma {
		cfg.GammaCoefficient = userCfg.Gamma
	}
	if cfg.SymmetryLevel == ZeroSymmetryLevel {
		cfg.SymmetryLevel = userCfg.SymmetryLevel
	}
	if cfg.Functions == ZeroFunctions {
		cfg.Functions = FunctionsToString(userCfg.Functions)
	}
	if cfg.AffineParams == ZeroAffineParams {
		cfg.AffineParams = AffinesToString(userCfg.AffineParams)
	}
	return cfg
}

func FunctionsToString(funks []dto.Function) string {
	var b strings.Builder

	for i, function := range funks {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(function.Name)
		b.WriteByte(':')
		b.WriteString(strconv.FormatFloat(function.Weight, 'f', -1, 64))
	}
	return b.String()
}

func AffinesToString(affines []dto.AffineParams) string {
	var b strings.Builder
	for i, affine := range affines {
		if i > 0 {
			b.WriteByte('/')
		}

		b.WriteString(strconv.FormatFloat(affine.A, 'f', -1, 64))
		b.WriteByte(',')
		b.WriteString(strconv.FormatFloat(affine.B, 'f', -1, 64))
		b.WriteByte(',')
		b.WriteString(strconv.FormatFloat(affine.C, 'f', -1, 64))
		b.WriteByte(',')
		b.WriteString(strconv.FormatFloat(affine.D, 'f', -1, 64))
		b.WriteByte(',')
		b.WriteString(strconv.FormatFloat(affine.E, 'f', -1, 64))
		b.WriteByte(',')
		b.WriteString(strconv.FormatFloat(affine.F, 'f', -1, 64))
	}
	return b.String()
}
