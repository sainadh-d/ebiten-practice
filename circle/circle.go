package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// https://golangprojectstructure.com/creating-cool-games-with-ebiten-in-go/

const (
	screenWidth  = 300
	screenHeight = 300
)

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) drawCircle(screen *ebiten.Image, x, y, radius int, clr color.Color) error {
	radius64 := float64(radius)
	minAngle := math.Acos(1 - 1/radius64)
	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius64 * math.Cos(angle)
		yDelta := radius64 * math.Sin(angle)
		x1 := int(math.Round(float64(x) + xDelta))
		y1 := int(math.Round(float64(y) + yDelta))
		screen.Set(x1, y1, clr)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	red := color.RGBA{255, 0, 0, 0}
	g.drawCircle(screen, 150, 150, 50, red)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ebiten Hello")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
