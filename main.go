package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth  int = 320
	screenHeight int = 240
)

var fontFace text.Face = text.NewGoXFace(bitmapfont.Face)

type World struct {
	grid   []bool
	width  int
	height int
}

type Game struct {
	world *World
	keys  []ebiten.Key
}

func NewWorld(width, height int) *World {
	w := &World{
		make([]bool, width*height),
		width,
		height,
	}

	for i := range w.grid {
		w.grid[i] = rand.Intn(100) < 40
	}

	return w
}

func (w *World) countNeighbors(i int) int {
	count := 0
	y := i / w.width
	x := i % w.width
	bounds := w.width * w.height

	// (x0, y0) = (-1, -1), (0, -1), (1, -1), ..., (1, 1)
	for y0 := -1; y0 <= 1; y0++ {
		for x0 := -1; x0 <= 1; x0++ {
			if x0 == 0 && y0 == 0 {
				continue
			}

			iNeighbor := (y+y0)*w.width + (x + x0)
			if iNeighbor <= 0 || iNeighbor >= bounds {
				continue
			}
			if w.grid[iNeighbor] {
				count++
			}
		}
	}

	return count
}

func (g *Game) Update() error {
	// update board based on game of life laws
	next := make([]bool, g.world.width*g.world.height)
	for i := range g.world.grid {
		c := g.world.countNeighbors(i)

		switch {
		case c < 2 || c > 3:
			next[i] = false
		case c == 2:
			next[i] = g.world.grid[i]
		case c == 3:
			next[i] = true
		}
	}
	g.world.grid = next

	// reset on 'R' keypress
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.world = NewWorld(screenWidth, screenHeight)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw white pixels
	pixels := make([]byte, 4*g.world.width*g.world.height)
	for i, v := range g.world.grid {
		if v {
			pixels[4*i] = 0xff
			pixels[4*i+1] = 0xff
			pixels[4*i+2] = 0xff
			pixels[4*i+3] = 0xff
		} else {
			pixels[4*i] = 0
			pixels[4*i+1] = 0
			pixels[4*i+2] = 0
			pixels[4*i+3] = 0
		}
	}
	screen.WritePixels(pixels)

	// text
	text.Draw(screen, "press R to reset world", fontFace, &text.DrawOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		NewWorld(screenWidth, screenHeight),
		[]ebiten.Key{},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("conway")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
