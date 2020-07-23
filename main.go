package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/szTheory/chip8go/emu"
)

const (
	scaleFactor    = 10
	cyclesPerFrame = 9
	// cyclesPerFrame = 60
)

type Game struct {
	emulator *emu.Emulator
	canvas   *ebiten.Image
}

func (g *Game) Setup(romFilename string) {
	g.emulator = new(emu.Emulator)
	g.emulator.Setup(romFilename)

	var err error
	if g.canvas, err = ebiten.NewImage(emu.ScreenWidthPx, emu.ScreenHeightPx, ebiten.FilterDefault); err != nil {
		panic(err)
	}
	if err := g.canvas.Fill(color.Black); err != nil {
		panic(err)
	}
}

// Update the logical state
func (g *Game) Update(screen *ebiten.Image) error {
	// start := time.Now()
	for i := 0; i < cyclesPerFrame; i++ {
		g.emulator.EmulateCycle()
	}
	// elapsed := time.Since(start)
	// fmt.Printf("Update time ms: %d\n", elapsed.Milliseconds())
	// fmt.Printf("TPS: %f", ebiten.CurrentTPS())

	return nil
}

// Render the screen
func (g *Game) Draw(screen *ebiten.Image) {
	// start := time.Now()
	for x := 0; x < emu.ScreenWidthPx; x++ {
		for y := 0; y < emu.ScreenHeightPx; y++ {
			setColor := color.Black
			if g.emulator.Display.Pixels[x][y] == 1 {
				setColor = color.White
			}
			if setColor != g.canvas.At(x, y) {
				g.canvas.Set(x, y, setColor)
			}
		}
	}

	geometry := ebiten.GeoM{}
	if err := screen.DrawImage(g.canvas, &ebiten.DrawImageOptions{GeoM: geometry}); err != nil {
		panic(err)
	}
	fmt.Printf("FPS: %f", ebiten.CurrentFPS())
	// elapsed := time.Since(start)
	// fmt.Printf("Draw time ms: %d\n", elapsed.Milliseconds())

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return emu.ScreenWidthPx, emu.ScreenHeightPx
}

const (
	ScreenWidthPx  = emu.ScreenWidthPx * scaleFactor
	ScreenHeightPx = emu.ScreenHeightPx * scaleFactor
)

func main() {
	// romFilename := "roms/PONG.ch8"
	// romFilename := "roms/test_opcode.ch8"
	// romFilename := "roms/BC_test.ch8"
	// romFilename := "roms/IBM.ch8"
	// romFilename := "roms/TETRIS.ch8"
	// romFilename := "roms/LANDING.ch8"
	romFilename := "roms/KALEID.ch8"

	ebiten.SetWindowSize(ScreenWidthPx, ScreenHeightPx)
	ebiten.SetWindowTitle("Chip-8 - " + romFilename)

	game := new(Game)
	game.Setup(romFilename)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
