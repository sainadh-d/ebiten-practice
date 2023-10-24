package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

	run bool
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

	// Go through all the cells
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			op := &ebiten.DrawImageOptions{}

			x := g.startX + r*g.cellSize
			y := g.startY + c*g.cellSize

			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(outer, op)

			cell := Cell{r, c}
			if found := g.liveCells[cell]; found {
				// Don't draw inner rectangle if this is a live cell.
				// This allows us to see the live cell as a filled rectangle with white color
			} else {
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Translate(float64(x+g.edgeWidth), float64(y+g.edgeWidth))
				screen.DrawImage(inner, op2)
			}
		}
	}
}

func (g *Grid) update() {
	if !g.run {
		return
	}

	timeDelta := time.Since(lastUpdatedTime)

	if timeDelta < time.Millisecond*200 {
		return
	} else {
		lastUpdatedTime = time.Now()
	}

	// @Speed Can we figure out a way not to make a new map everytime?
	nextGen := make(map[Cell]bool)

	// Go through all the cells
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			neighbours := getNeighbours(Cell{r, c})

			liveNeighborCount := 0
			for i := range neighbours {
				if g.liveCells[neighbours[i]] {
					liveNeighborCount++
				}
			}

			// Apply the rules
			isAlive := g.liveCells[Cell{r, c}]
			if isAlive && liveNeighborCount == 2 || liveNeighborCount == 3 {
				// Cell continues to stay alive don't do anything
				nextGen[Cell{r, c}] = true
			} else if !isAlive && liveNeighborCount == 3 {
				// Cell becomes alive
				nextGen[Cell{r, c}] = true
			} else {
				// Cell becomes dead
			}
		}
	}
	g.liveCells = nextGen

}

func (g *Grid) handleKeyEvent(key ebiten.Key) {
	switch {
	case key == ebiten.KeyC:
		// Clear the map
		g.liveCells = make(map[Cell]bool)
	case key == ebiten.KeySpace:
		g.run = !g.run
	}
}

func (g *Grid) handleMouseEvent(mx, my int) {
	x := (mx - g.startX) / g.cellSize
	y := (my - g.startY) / g.cellSize

	if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
		// Out of grid area - Do nothing
		return
	}
	c := Cell{x, y}
	g.liveCells[c] = !g.liveCells[c]
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
		startY: 80,

		rows: 30,
		cols: 30,

		cellSize:  20,
		edgeWidth: 1,

		liveCells: make(map[Cell]bool),
	}

	TechnoRaceSmall  font.Face
	TechnoRaceNormal font.Face
	TechnoRaceBig    font.Face
)

// ----------------- Init --------------------
func init() {
	fontData, err := os.ReadFile("techno-race.otf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	TechnoRaceSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	TechnoRaceNormal, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	TechnoRaceBig, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull, // Use quantization to save glyph cache images.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adjust the line height.
	TechnoRaceBig = text.FaceWithLineHeight(TechnoRaceBig, 32)
}

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

	if repeatingButtonPressed(ebiten.MouseButtonLeft) {
		grid.handleMouseEvent(mx, my)
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyC) {
		grid.handleKeyEvent(ebiten.KeyC)
		return nil
	}

	if repeatingKeyPressed(ebiten.KeySpace) {
		grid.handleKeyEvent(ebiten.KeySpace)
	}

	// @Cleanup: The below code should also probably be moved to grid.update()

	grid.update()

	return nil
}

// https://github.com/sedyh/ebitengine-cheatsheet#center-text
func DrawCenteredText(screen *ebiten.Image, textFont font.Face, s string, cx, cy int) {
	// @Cleanup
	/*
				bounds, _ := font.BoundString(textFont, s)
				dx := (bounds.Max.X - bounds.Min.X)
				dy := (bounds.Max.Y - bounds.Min.Y)

				x, y := fixed.Int26_6(cx)-bounds.Min.X-dx/2, fixed.Int26_6(cy)-bounds.Min.Y-dy/2
		        text.Draw(screen, s, textFont, int(x), int(y), colornames.Red)
	*/
	bounds := text.BoundString(textFont, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, textFont, x, y, colornames.White)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background Color
	// screen.Fill(color.RGBA{255, 255, 255, 0})

	DrawCenteredText(screen, TechnoRaceBig, "GAME OF LIFE", screenWidth/2, 20)

	msg := "Press Space to START or STOP, C to CLEAR"
	DrawCenteredText(screen, TechnoRaceNormal, msg, screenWidth/2, 50)

	// Draw Status
	if grid.run {
		msg = "Status:  Running"
	} else {
		msg = "Status:  Stopped"
	}
	bounds := text.BoundString(TechnoRaceSmall, msg)
	text.Draw(screen, msg, TechnoRaceSmall, screenWidth-bounds.Dx()-20, 20, color.White)

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
