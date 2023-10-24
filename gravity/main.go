package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1346
	screenHeight = 768
	blockSize    = 3

	xAcc = 0
	yAcc = float64(250)

	e = 0.8 // Coefficient of Restiution
)

var (
	X float64 = 100
	Y float64 = 100

	vx float64 = 0
	vy float64 = 0

	img *ebiten.Image

	lastUpdatedTime = time.Now()
)

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("football.png")
	fmt.Println(fmt.Sprintf("Image Size: X: %d Y: %d", img.Bounds().Dx(), img.Bounds().Dy()))
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

	// ------ Input Stuff
	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "ArrowDown":
			// @Incomplete: Y Should be min(screenHeight-ballDiameter, Y+1). But I don't know
			// how to figure out ball radius from image yet.
			Y = min(screenHeight, Y+1)
		case "ArrowUp":
			Y = max(0, Y-1)
		case "ArrowRight":
			// @Incomplete: Y Should be min(screenHeight-ballDiameter, Y+1). But I don't know
			// how to figure out ball radius from image yet.
			X = min(screenWidth-blockSize, X+1)
		case "ArrowLeft":
			X = max(0, X-1)
		}
	}

	timeDelta := float64(time.Since(lastUpdatedTime))
	timeDelta = timeDelta / 1000000000

	lastUpdatedTime = time.Now()

	// Gravity - Kinematics.

	vx = vx + xAcc*timeDelta
	vy = vy + yAcc*timeDelta

	X += vx * timeDelta
	Y += vy * timeDelta

	// @Incomplete implement in all directions
	if Y >= screenHeight-64 {
		// We touched the ground. Now we jump
		vy = -e * vy
		Y = screenHeight - 64
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
	// fmt.Printf("Screen Size X: %d, Y: %d\n", screen.Bounds().Dx(), screen.Bounds().Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(X, Y)
	// op.ColorScale.Scale(0.5, 0.5, 0.5, 1)

	screen.DrawImage(img, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%.2f, y: %.2f, FPS: %.2f)", X, Y, ebiten.ActualFPS()))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Gravity")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
