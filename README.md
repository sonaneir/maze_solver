# maze_solver

A command-line program in Go that solves mazes given as images. It reads a maze from a PNG file, finds a path from the entrance to the exit, and saves a new image with the solution path drawn on it.

## How it works

The maze is a PNG image where walls, open paths, and the start/end points are represented by different colors. The solver reads the image, treats it as a grid, and searches for a path from start to finish using a pathfinding algorithm over the grid. Once a path is found, it draws it onto a copy of the image and writes the result to an output file.

The repository includes example solutions — `sol.png` (a solved maze) and `sol.gif` (an animation of the solving process).

## Usage

```bash
maze_solver input.png output.png
```

- `input.png` — the maze image to solve
- `output.png` — where the solved maze image is written

If the arguments are missing, the program prints usage and exits.

Example:

```bash
go run . mazes/maze.png solution.png
```

## Project structure

```
maze_solver/
├── main.go                 # CLI entry point: parses args, runs the solver, saves the result
├── internal/
│   └── solver/             # core logic: load image, find path, save solution
└── mazes/                  # example maze images
```

The solving logic is isolated in `internal/solver` behind a small API — `New` to load a maze, `Solve` to find the path, and `SaveSolution` to write the output image — so the CLI in `main.go` stays thin.
