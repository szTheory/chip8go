package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/szTheory/chip8go/emu"
)

const (
	scaleFactor    = 10
	cyclesPerFrame = 9
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

// Update the logical state
func (g *Game) Update(screen *ebiten.Image) error {
	inputs := keyPairs()
	for i := 0; i < len(inputs); i++ {
		keyIndex := inputs[i].index
		key := inputs[i].key
		isPressed := ebiten.IsKeyPressed(key)

		g.emulator.Input.Update(keyIndex, isPressed)
		if isPressed && inpututil.IsKeyJustPressed(key) && g.emulator.Input.WaitingForInput {
			g.emulator.CatchInput(keyIndex)
		}
	}

	for i := 0; i < cyclesPerFrame; i++ {
		g.emulator.EmulateCycle()
	}

	var volume float64
	if g.emulator.SoundEnabled() {
		volume = 1
	}
	g.emulator.AudioPlayer.SetVolume(volume)

	g.emulator.UpdateDelayTimer()
	g.emulator.UpdateSoundTimer()

	return nil
}

// Render the screen
func (g *Game) Draw(screen *ebiten.Image) {
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return emu.ScreenWidthPx, emu.ScreenHeightPx
}

const (
	ScreenWidth  = emu.ScreenWidthPx * scaleFactor
	ScreenHeight = emu.ScreenHeightPx * scaleFactor
)

func main() {
	// romFilename := "roms/PONG.ch8"
	// romFilename := "roms/test_opcode.ch8"
	// romFilename := "roms/BC_test.ch8"
	// romFilename := "roms/IBM.ch8"
	romFilename := "roms/TETRIS.ch8"
	// romFilename := "roms/LANDING.ch8"
	// romFilename := "roms/KALEID.ch8"
	// romFilename := "roms/TRON.ch8"

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Chip-8 - " + romFilename)

	game := new(Game)
	game.Setup(romFilename)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
