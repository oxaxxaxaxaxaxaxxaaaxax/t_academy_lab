package dto

type UserConfig struct {
	Size               Size           `json:"size"`
	IterationCount     int            `json:"iteration_count"`
	OutputPath         string         `json:"output_path"`
	Threads            int            `json:"threads"`
	Seed               int            `json:"seed"`
	Functions          []Function     `json:"functions"`
	AffineParams       []AffineParams `json:"affine_params"`
	Gamma              float64        `json:"gamma"`
	AddGammaCorrection bool           `json:"gamma_correction"`
	SymmetryLevel      int            `json:"symmetry_level"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Function struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type AffineParams struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
	D float64 `json:"d"`
	E float64 `json:"e"`
	F float64 `json:"f"`
}
