package domain

import "sync"

type Pixel struct {
	HitCounter      int
	NormalizedValue float64
	R, G, B         uint8
}

type Canvas struct {
	Pixels  Pixels
	Mutexes []sync.Mutex
}
type Pixels [][]Pixel

func NewCanvas(width, height int) Canvas {
	canvas := Canvas{
		Pixels:  make([][]Pixel, height),
		Mutexes: make([]sync.Mutex, height),
	}
	for i := range canvas.Pixels {
		canvas.Pixels[i] = make([]Pixel, width)
	}
	return canvas
}
