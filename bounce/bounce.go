package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// https://golangprojectstructure.com/creating-cool-games-with-ebiten-in-go/

const (
	screenWidth  = 300
	screenHeight = 300
	ballRadius   = 15
)

var (
	ballPositionX = float64(screenWidth) / 2
	ballPositionY = float64(screenHeight) / 2

	ballMovementX = float64(0.00000006)
	ballMovementY = float64(0.00000004)

	lastUpdatedTime = time.Now()
)

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	timeDelta := float64(time.Since(lastUpdatedTime))
	lastUpdatedTime = time.Now()

	ballPositionX += ballMovementX * timeDelta
	ballPositionY += ballMovementY * timeDelta

	const minX = ballRadius
	const minY = ballRadius
	const maxX = screenWidth - ballRadius
	const maxY = screenHeight - ballRadius

	if ballPositionX >= maxX || ballPositionX <= minX {
		if ballPositionX > maxX {
			ballPositionX = maxX
		} else if ballPositionX < minX {
			ballPositionX = minX
		}
		ballMovementX *= -1
	}

	if ballPositionY >= maxY || ballPositionY <= minY {
		if ballPositionY > maxY {
			ballPositionY = maxY
		} else if ballPositionY < minY {
			ballPositionY = minY
		}

		ballMovementY *= -1
	}
	return nil
}

// drawCircle - Draws a Circle with x, y as Center
func (g *Game) drawCircle(screen *ebiten.Image, x, y, radius int, clr color.Color) error {
	radius64 := float64(radius)
	minAngle := math.Acos(1 - 1/radius64)
	white := color.RGBA{255, 255, 255, 0}
	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius64 * math.Cos(angle)
		yDelta := radius64 * math.Sin(angle)
		x1 := int(math.Round(float64(x) + xDelta))
		y1 := int(math.Round(float64(y) + yDelta))

		// Fill the circle
		if y1 < y {
			for y2 := y1; y2 <= y+(y-y1); y2++ {
				screen.Set(x1, y2, clr)
			}
		} else {
			/*
				for y2 := y1; y2 > y; y2-- {
					screen.Set(x1, y2, clr)
				}
			*/
		}

		screen.Set(x1, y1, clr)
	}
	screen.Set(x, y, white)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	red := color.RGBA{255, 0, 0, 0}
	x := int(math.Round(ballPositionX))
	y := int(math.Round(ballPositionY))
	g.drawCircle(screen, x, y, ballRadius, red)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ebiten Hello")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
