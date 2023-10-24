package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	blockSize    = 3
)

var (
	X float64 = 100
	Y float64 = 100

	img *ebiten.Image
)

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}
}

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
			Y = min(screenHeight-blockSize, Y+1)
		case "ArrowUp":
			Y = max(0, Y-1)
		case "ArrowRight":
			X = min(screenWidth-blockSize, X+1)
		case "ArrowLeft":
			X = max(0, X-1)
		}
	}

	return nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 0})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(X, Y)
	// op.ColorScale.Scale(0.5, 0.5, 0.5, 1)

	screen.DrawImage(img, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%.2f, y: %.2f, FPS: %.2f)", X, Y, ebiten.ActualFPS()))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
