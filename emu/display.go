package emu

const (
	ScreenWidthPx  = 64
	ScreenHeightPx = 32

	SpriteWidthPx = 8
)

type Display struct {
	Pixels [ScreenWidthPx][ScreenHeightPx]byte
	Draw   bool
}

func (d *Display) Clear() {
	for x := 0; x < ScreenWidthPx; x++ {
		for y := 0; y < ScreenHeightPx; y++ {
			d.Pixels[x][y] = 0
		}
	}
}

// Sprites are XORed onto the existing screen.
// Returns true if any pixels were erased, false otherwise
func (d *Display) DrawSprite(x byte, y byte, row byte) bool {
	erased := false

	for i := x; i < x+8; i++ {
		xIndex := i % ScreenWidthPx
		yIndex := y % ScreenHeightPx
		wasSet := d.Pixels[xIndex][yIndex] == 0x1
		value := byte(row >> (x + 8 - i - 1) & 0x1)

		d.Pixels[xIndex][yIndex] ^= value
		// newValue := d.Pixels[xIndex][y]
		// fmt.Println(newValue)

		if d.Pixels[xIndex][yIndex] == 0x0 && wasSet {
			erased = true
		}
	}

	return erased
}

// func (d *Display) Draw(mem *Memory) {
// 	if !mem.ShouldDraw() {
// 		return
// 	}
// }
