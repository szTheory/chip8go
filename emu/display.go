package emu

type Display struct {
	pixels [ScreenWidthPx][ScreenHeightPx]byte
}

const (
	ScreenWidthPx  = 64
	ScreenHeightPx = 32

	SpriteWidthPx = 8
	// SpriteHeightPx = 8
)

// Returns true if any set pixels were changed to unset
// false otherwise
func (d *Display) DrawSprite(x byte, y byte, row byte) bool {
	unset := false

	for i := x; i < 8; i++ {
		xIndex := x % ScreenWidthPx
		wasSet := d.pixels[xIndex][y] == 0x0001
		set := row >> (8 - i) & 0x0001
		d.pixels[xIndex][y] = set
		if set == 0x0001 && wasSet {
			unset = true
		}
	}

	return unset
}

// func (d *Display) Draw(mem *Memory) {
// 	if !mem.ShouldDraw() {
// 		return
// 	}
// }
