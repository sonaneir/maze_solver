package solver

import (
	"image"
	"log"
	"sync"
)

// listenToBranches creates a new goroutine for each branch published in
// s.pathsToExplore.
func (s *Solver) listenToBranches() {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for {
		select {
		case <-s.quit:
			log.Println("the treasure has been found, stopping worker")
			return
		case p := <-s.pathsToExplore:
			wg.Add(1)
			go func(path *path) {
				defer wg.Done()

				s.explore(p)
			}(p)
		}
	}
}

// explore one path and publish to the s.pathsToExplore channel
// any branch we discover that we don't take.
func (s *Solver) explore(pathToBranch *path) {
	if pathToBranch == nil {
		return
	}

	pos := pathToBranch.at

	for {
		select {
		case <-s.quit:
			return
		case s.exploredPixels <- pos:
		}

		candidates := make([]image.Point, 0, 3)
		for _, n := range neighbours(pos) {
			if pathToBranch.isPreviousStep(n) {
				continue
			}

			switch s.maze.RGBAAt(n.X, n.Y) {
			case s.palette.treasure:
				s.mutex.Lock()
				if s.solution == nil {
					s.solution = &path{previousStep: pathToBranch, at: n}
					log.Printf("Treasure found: %v!", s.solution.at)
					close(s.quit)
				}
				s.mutex.Unlock()
				return
			case s.palette.path:
				candidates = append(candidates, n)
			}
		}

		if len(candidates) == 0 {
			log.Printf("I must have taken the wrong turn at position %v.", pos)
			return
		}

		for _, candidate := range candidates[1:] {
			branch := &path{previousStep: pathToBranch, at: candidate}
			// We are sure we send to pathsToExplore only when the quit channel isn't closed.
			// A goroutine might have found the treasure since the check at the start of the loop.
			select {
			// s.quit returns a zero value only when the channel was closed, here -- when the exploration should end.
			case <-s.quit:
				log.Printf("I'm an unlucky branch, someone else found the treasure, I give up at position %v.", pos)
				return
			case s.pathsToExplore <- branch:
				// continue execution after the select block
			}
		}

		pathToBranch = &path{previousStep: pathToBranch, at: candidates[0]}
		pos = candidates[0]
	}
}
