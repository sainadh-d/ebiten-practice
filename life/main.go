package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// @Speed: Don't use a struct for Cell. For index should suffice and
// x, y can be inferred based on the Cell index
type Cell struct {
	x int
	y int
}

type Grid struct {
	startX int
	startY int
	rows   int
	cols   int

	cellSize  int
	edgeWidth int

	// @Speed: Can we just use a bitmap?
	liveCells map[Cell]bool
}

func (g *Grid) draw(screen *ebiten.Image) {
	// Outer
	// (0, 0), (75, 0), (150, 0), (225, 0)
	// (0, 75), (75, 75), (150, 75), (225, 75)

	// x * w
	// y * w

	// Inner
	// (1, 1), (76, 1), (151, 1), (225, 0)
	// (1, 76), (75, 75), (150, 75), (225, 75)

	// 1 + x * w
	// 1 + y * w

	// Draw the Grid

	outer := ebiten.NewImage(g.cellSize, g.cellSize)
	outer.Fill(color.RGBA{101, 107, 117, 0})
	inner := ebiten.NewImage(g.cellSize-g.edgeWidth-1, g.cellSize-g.edgeWidth-1)
	inner.Fill(color.Black)

	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			op := &ebiten.DrawImageOptions{}

			x := g.startX + r*g.cellSize
			y := g.startY + c*g.cellSize

			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(outer, op)

			c := Cell{x, y}
			if found := g.liveCells[c]; found {
				// Don't draw inner rectangle if this is a live cell.
				// This allows us to see the live cell as a filled rectangle with white color
			} else {
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Translate(float64(x+g.edgeWidth), float64(y+g.edgeWidth))
				screen.DrawImage(inner, op2)
			}
		}
	}

	/*
		for x := 0; x < (screenWidth / g.cellSize); x += 1 {
			for y := 0; y < (screenHeight / g.cellSize); y += 1 {
				c := Cell{x, y}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*g.cellSize), float64(y*g.cellSize))
				screen.DrawImage(outer, op)

				if found := g.liveCells[c]; found {
					// Don't draw inner rectangle if this is a live cell.
					// This allows us to see the live cell as a filled rectangle with white color
				} else {
					op2 := &ebiten.DrawImageOptions{}
					op2.GeoM.Translate(float64(x*g.cellSize+1), float64(y*g.cellSize+1))
					screen.DrawImage(inner, op2)
				}
			}
		}

	*/
}

func (g *Grid) update() {
	// @Speed Can we figure out a way not to make a new map everytime?
	nextGen := make(map[Cell]bool)
	// Go through all the cells
	for x := 0; x < (screenWidth / g.cellSize); x += 1 {
		for y := 0; y < (screenHeight / g.cellSize); y += 1 {
			// Get the neighbours of this cell
			neighbours := getNeighbours(Cell{x, y})

			liveNeighborCount := 0
			for i := range neighbours {
				if g.liveCells[neighbours[i]] {
					liveNeighborCount++
				}
			}

			// Apply the rules
			isAlive := g.liveCells[Cell{x, y}]
			if isAlive && liveNeighborCount == 2 || liveNeighborCount == 3 {
				// Cell continues to stay alive don't do anything
				nextGen[Cell{x, y}] = true
			} else if !isAlive && liveNeighborCount == 3 {
				// Cell becomes alive
				nextGen[Cell{x, y}] = true
			} else {
				// Cell becomes dead
			}
		}
	}
	g.liveCells = nextGen
}

// ---------------- Variables --------------------
const (
	screenWidth  = 720
	screenHeight = 720
)

var (
	start = false

	lastUpdatedTime = time.Now()

	grid = &Grid{
		startX: 60,
		startY: 60,

		rows: 30,
		cols: 30,

		cellSize:  20,
		edgeWidth: 1,

		liveCells: make(map[Cell]bool),
	}
)

// ------------- Utils -------------------------

func getNeighbours(c Cell) []Cell {
	neighbours := []Cell{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			neighbours = append(neighbours, Cell{c.x + x, c.y + y})
		}
	}
	return neighbours
}

// repeatingButtonPressed return true when key is pressed considering the repeat state.
func repeatingButtonPressed(button ebiten.MouseButton) bool {
	const (
		delay    = 30
		interval = 3
		Grid
	)
	d := inpututil.MouseButtonPressDuration(button)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

// ----------------- Game -------------------------

type Game struct {
	pressedKeys []ebiten.Key
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	/*
		g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

		for _, key := range g.pressedKeys {
			switch key.String() {
			case "Space":
				start = !start
			case "ArrowUp":
			case "ArrowRight":
			case "ArrowLeft":
			}
		}
	*/

	// @Cleanup
	mx, my := ebiten.CursorPosition()

	// Grid Coords
	cx := mx / grid.cellSize
	cy := my / grid.cellSize

	if repeatingButtonPressed(ebiten.MouseButtonLeft) {
		cell := Cell{cx, cy}
		grid.liveCells[cell] = !grid.liveCells[cell]
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyC) {
		// Clear the map
		grid.liveCells = map[Cell]bool{}
		return nil
	}

	if repeatingKeyPressed(ebiten.KeySpace) {
		fmt.Println("Pressed Space, toggling to", start)
		start = !start
	}

	if !start {
		return nil
	}

	timeDelta := time.Since(lastUpdatedTime)

	if timeDelta < time.Millisecond*200 {
		return nil
	} else {
		lastUpdatedTime = time.Now()
	}

	grid.update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background Color
	// screen.Fill(color.RGBA{255, 255, 255, 0})

	// Draw the grid
	grid.draw(screen)

}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game of Life")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
