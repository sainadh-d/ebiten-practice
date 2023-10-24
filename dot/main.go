package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	blockSize = 3
)

var (
	x = 100
	y = 100
)

type Game struct {
	pressedKeys []ebiten.Key
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {

	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "ArrowDown":
			y = min(screenHeight-blockSize, y+1)
		case "ArrowUp":
			y = max(0, y-1)
		case "ArrowRight":
			x = min(screenWidth-blockSize, x+1)
		case "ArrowLeft":
			x = max(0, x-1)
		}
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 0})

	// Draw Rectangle with x and y as center
	for x1 := -blockSize; x1 <= blockSize; x1++ {
		for y1 := -blockSize; y1 <= blockSize; y1++ {
			if x1 == -blockSize || x1 == blockSize || y1 == -blockSize || y1 == blockSize {
				screen.Set(x+x1, y+y1, color.RGBA{0, 0, 255, 0})
			}
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, y: %d, FPS: %.2f)", x, y, ebiten.ActualFPS()))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
