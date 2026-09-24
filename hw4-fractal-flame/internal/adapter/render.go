package adapter

import (
	"hw4-fractal-flame/internal/domain"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
)

type Render struct {
	Image image.NRGBA
}

func NewRender(cfg domain.FractalConfig) Render {
	var render Render

	render.Image = *image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
	slog.Debug("render conf", "width", cfg.Width, "height", cfg.Height)
	return render
}

func (r Render) RenderImage(pixels domain.Pixels, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	for row, p := range pixels {
		for col, pix := range p {
			r.Image.Set(col, row, color.NRGBA{R: pix.R, G: pix.G, B: pix.B, A: 255})
		}
	}
	err = png.Encode(f, &r.Image)
	if err != nil {
		return err
	}
	return nil
}
