package domain

type FractalConfig struct {
	Width             int
	Height            int
	Seed              int64
	IterCount         int
	OutputPath        string
	Threads           int
	AffineParams      string
	Functions         string
	ConfigPath        string
	IterAlgorithmName string
	GammaCoefficient  float64
	SymmetryLevel     int
}

func NewFractalConfig(width, height, iterCount, threads, symmetry int, seed int64,
	output, affine, functions, configPath, iterAlgorithmName string, gamma float64) FractalConfig {
	return FractalConfig{
		Width:             width,
		Height:            height,
		Seed:              seed,
		IterCount:         iterCount,
		Threads:           threads,
		OutputPath:        output,
		AffineParams:      affine,
		Functions:         functions,
		ConfigPath:        configPath,
		IterAlgorithmName: iterAlgorithmName,
		GammaCoefficient:  gamma,
		SymmetryLevel:     symmetry,
	}
}
