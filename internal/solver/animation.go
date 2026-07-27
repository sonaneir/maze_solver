package solver

import (
	"image"
	plt "image/color/palette"

	"golang.org/x/image/draw"
)

// countExplorablePixels scans the maze and counts the number
// of pixels that are not walls.
func (s *Solver) countExplorablePixels() int {
	explorablePixels := 0
	for row := s.maze.Bounds().Min.Y; row < s.maze.Bounds().Max.Y; row++ {
		for col := s.maze.Bounds().Min.X; col < s.maze.Bounds().Max.X; col++ {
			if s.maze.RGBAAt(col, row) != s.palette.wall {
				explorablePixels++
			}
		}
	}

	return explorablePixels
}

// registerExploredPixels registers positions as explored on the image,
// and, if we reach a threshold, adds the frame to the output GIF.
func (s *Solver) registerExploredPixels() {
	const totalExpectedFrames = 30

	explorablePixels := s.countExplorablePixels()
	pixelsExplored := 0

	frameInterval := explorablePixels / totalExpectedFrames
	if frameInterval == 0 {
		frameInterval = 1
	}

	for {
		select {
		case <-s.quit:
			return
		case pos := <-s.exploredPixels:
			s.maze.Set(pos.X, pos.Y, s.palette.explored)
			pixelsExplored++
			if pixelsExplored%frameInterval == 0 {
				s.drawCurrentFrameToGIF()
			}
		}
	}
}

// drawCurrentFrameToGIF adds the current state of the maze as a frame of the animation.
func (s *Solver) drawCurrentFrameToGIF() {
	const (
		gifWidth      = 500
		frameDuration = 20
	)

	frame := image.NewPaletted(image.Rect(0, 0, gifWidth, gifWidth*s.maze.Bounds().Dy()/s.maze.Bounds().Dx()), plt.Plan9)

	draw.NearestNeighbor.Scale(frame, frame.Rect, s.maze, s.maze.Bounds(), draw.Over, nil)

	s.animation.Image = append(s.animation.Image, frame)
	s.animation.Delay = append(s.animation.Delay, frameDuration)
}

// writeLastFrame writes the last frame of the gif, with the solution highlighted.
func (s *Solver) writeLastFrame() {
	stepsFromTreasure := s.solution
	// Paint the path from entrance to the treasure.
	for stepsFromTreasure != nil {
		s.maze.Set(stepsFromTreasure.at.X, stepsFromTreasure.at.Y, s.palette.solution)
		stepsFromTreasure = stepsFromTreasure.previousStep
	}

	const solutionFrameDuration = 300 // 3 seconds
	// Add the solution frame, with the coloured path, to the output gif.
	s.drawCurrentFrameToGIF()
	s.animation.Delay[len(s.animation.Delay)-1] = solutionFrameDuration
}
