package main

import (
	"errors"
	"image/color"
	"path"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/sqweek/dialog"
	"github.com/szTheory/chip8go/emu"
)

func main() {
	game := new(Game)
	if err := game.pickGame(); err != nil {
		return
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

const (
	scaleFactor    = 10
	cyclesPerFrame = 10
	ScreenWidth    = emu.ScreenWidthPx * scaleFactor
	ScreenHeight   = emu.ScreenHeightPx * scaleFactor
)

type Game struct {
	emulator    *emu.Emulator
	romFilename string
}

// Update the logical state
func (g *Game) Update(screen *ebiten.Image) error {
	inputs := keyPairs()

	// Enter key resets game
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.reset()
	}

	for i := 0; i < cyclesPerFrame; i++ {
		// update inputs
		for i := 0; i < len(inputs); i++ {
			keyIndex := inputs[i].index
			key := inputs[i].key
			isPressed := ebiten.IsKeyPressed(key)

			g.emulator.Input.Update(keyIndex, isPressed)
			if isPressed && inpututil.IsKeyJustPressed(key) && g.emulator.Input.WaitingForInput {
				g.emulator.CatchInput(keyIndex)
			}
		}

		// emulate a cycle
		g.emulator.EmulateCycle()
	}

	// update audio
	var volume float64
	if g.emulator.SoundEnabled() {
		volume = 1
	}
	g.emulator.AudioPlayer.SetVolume(volume)

	// update timers
	g.emulator.UpdateDelayTimer()
	g.emulator.UpdateSoundTimer()

	return nil
}

// Render the screen
func (g *Game) Draw(screen *ebiten.Image) {
	var err error
	var canvas *ebiten.Image
	if canvas, err = ebiten.NewImage(emu.ScreenWidthPx, emu.ScreenHeightPx, ebiten.FilterDefault); err != nil {
		panic(err)
	}
	if err := canvas.Fill(color.Black); err != nil {
		panic(err)
	}

	for x := 0; x < emu.ScreenWidthPx; x++ {
		for y := 0; y < emu.ScreenHeightPx; y++ {
			setColor := color.Black
			if g.emulator.Display.Pixels[x][y] == 1 {
				setColor = color.White
			}
			if setColor != canvas.At(x, y) {
				canvas.Set(x, y, setColor)
			}
		}
	}

	geometry := ebiten.GeoM{}
	if err := screen.DrawImage(canvas, &ebiten.DrawImageOptions{GeoM: geometry}); err != nil {
		panic(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return emu.ScreenWidthPx, emu.ScreenHeightPx
}

func (g *Game) reset() {
	g.emulator = new(emu.Emulator)
	g.emulator.Setup(g.romFilename)
}

func (g *Game) pickGame() error {
	romFilename, err := dialog.File().Filter("CHIP-8 game file", "ch8").Load()
	if err != nil {
		return err
	}
	if romFilename == "" {
		return errors.New("No game selected")
	}

	g.loadGame(romFilename)
	return nil
}

func (g *Game) loadGame(romFilename string) {
	ebiten.SetWindowTitle("Chip-8 - " + path.Base(romFilename))
	g.romFilename = romFilename
	g.reset()
}

type keyPair struct {
	index byte
	key   ebiten.Key
}

func keyPairs() [16]keyPair {
	list := [16]keyPair{
		{0, ebiten.KeyX},
		{1, ebiten.Key1},
		{2, ebiten.Key2},
		{3, ebiten.Key3},
		{4, ebiten.KeyQ},
		{5, ebiten.KeyW},
		{6, ebiten.KeyE},
		{7, ebiten.KeyA},
		{8, ebiten.KeyS},
		{9, ebiten.KeyD},
		{0xA, ebiten.KeyZ},
		{0xB, ebiten.KeyC},
		{0xC, ebiten.Key4},
		{0xD, ebiten.KeyR},
		{0xE, ebiten.KeyF},
		{0xF, ebiten.KeyV},
	}

	return list
}
