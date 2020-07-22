package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/szTheory/chip8go/emu"
)

type Game struct {
	emulator *emu.Emulator
	canvas   *ebiten.Image
	// frameBuffer [emu.ScreenWidthPx * emu.ScreenHeightPx * 4]byte
}

func (g *Game) Setup(romFilename string) {
	g.emulator = new(emu.Emulator)
	g.emulator.Setup(romFilename)

	var err error
	if g.canvas, err = ebiten.NewImage(emu.ScreenWidthPx, emu.ScreenHeightPx, ebiten.FilterDefault); err != nil {
		panic(err)
	}
}

// Update the logical state
func (g *Game) Update(screen *ebiten.Image) error {
	g.emulator.EmulateCycle()

	return nil
}

// Render the screen
func (g *Game) Draw(screen *ebiten.Image) {
	if err := g.canvas.Fill(color.Black); err != nil {
		panic(err)
	}
	for x := 0; x < emu.ScreenWidthPx; x++ {
		for y := 0; y < emu.ScreenHeightPx; y++ {
			if g.emulator.Display.Pixels[x][y] == 1 {
				g.canvas.Set(x, y, color.White)
			}
		}
	}

	geometry := ebiten.GeoM{}
	geometry.Scale(10, 10)
	if err := screen.DrawImage(g.canvas, &ebiten.DrawImageOptions{GeoM: geometry}); err != nil {
		panic(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

const (
	ScreenWidth  = emu.ScreenWidthPx * 10
	ScreenHeight = emu.ScreenHeightPx * 10
)

func main() {
	romFilename := "roms/PONG.ch8"
	// romFilename := "roms/BC_test.ch8"

	ebiten.SetWindowSize(emu.ScreenWidthPx*10, emu.ScreenHeightPx*10)
	ebiten.SetWindowTitle("Chip-8 - " + romFilename)

	game := new(Game)
	game.Setup(romFilename)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
