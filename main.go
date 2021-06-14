package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mmeinzer/chipate/interpret"
)

var (
	width  = 64
	height = 32
	scale  = 10

	black      = uint8(0)
	white      = uint8(255)
	blackColor = color.Gray{Y: black}
	whiteColor = color.Gray{Y: white}
)

// Game implements ebiten.Game interface.
type Game struct {
	count int
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	pixels := make([]byte, width*height)
	for i := range pixels {
		if i%2 == 0 {
			pixels[i] = black
		} else {
			pixels[i] = white
		}
	}

	img2 := &image.Gray{Pix: interpret.ChipPixels, Stride: width, Rect: image.Rect(0, 0, width, height)}

	img := ebiten.NewImageFromImage(img2)

	screen.DrawImage(img, &ebiten.DrawImageOptions{})
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	go interpret.Run("ibm.ch8")
	game := &Game{}

	ebiten.SetWindowSize(width*scale, height*scale)
	ebiten.SetWindowTitle("Chip Ate")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
