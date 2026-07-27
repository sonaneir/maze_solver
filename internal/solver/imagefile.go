package solver

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
	"strings"
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

// SaveSolution saves the image as a PNG file with the solution path
// highlighted.
func (s *Solver) SaveSolution(outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create output image file at %s", outputPath)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("unable to close file: %w", closeErr))
		}
	}()

	stepsFromTresure := s.solution

	for stepsFromTresure != nil {
		s.maze.Set(stepsFromTresure.at.X, stepsFromTresure.at.Y, s.palette.solution)
		stepsFromTresure = stepsFromTresure.previousStep
	}

	err = png.Encode(f, s.maze)
	if err != nil {
		return fmt.Errorf("unable to write output image at %s: %w", outputPath, err)
	}

	gifPath := strings.Replace(outputPath, "png", "gif", -1)
	err = s.saveAnimation(gifPath)

	if err != nil {
		return fmt.Errorf("unable to write output animation at %s", gifPath)
	}

	return nil
}

// saveAnimation writes the gif file.
func (s *Solver) saveAnimation(gifPath string) error {
	outputImage, err := os.Create(gifPath)
	if err != nil {
		return fmt.Errorf("unable to create output gif at %s: %w", gifPath, err)
	}

	defer func() {
		if closeErr := outputImage.Close(); closeErr != nil {
			// Return err and closeErr, in worst case scenario.
			err = errors.Join(err, fmt.Errorf("unable to close file: %w", closeErr))
		}
	}()

	log.Printf("animation contains %d frames\n", len(s.animation.Image))
	err = gif.EncodeAll(outputImage, s.animation)
	if err != nil {
		return fmt.Errorf("unable to encode gif: %w", err)
	}

	return nil
}
