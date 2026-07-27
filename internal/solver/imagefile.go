package solver

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

// openMaze opens a RGBA png image from a path.
func openMaze(imagePath string) (*image.RGBA, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open image %s: %w", imagePath, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("unable to load input image from %s: %w", imagePath, err)
	}

	rgbaImage, ok := img.(*image.RGBA)
	if !ok {
		b := img.Bounds()
		rgbaImage = image.NewRGBA(b)
		draw.Draw(rgbaImage, b, img, b.Min, draw.Src)
	}
	return rgbaImage, nil
}

// SaveSolution saves the image as a PNG file
// with the solution path highlighted.
func (s *Solver) SaveSolution(outputPath string) error {
	return nil
}
