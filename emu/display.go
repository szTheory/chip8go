package emu

type Display struct {
	Pixels [ScreenWidthPx][ScreenHeightPx]byte
}

const (
	ScreenWidthPx  = 64
	ScreenHeightPx = 32

	SpriteWidthPx = 8
)

func (d *Display) Clear() {
	for x := 0; x < ScreenWidthPx; x++ {
		for y := 0; y < ScreenHeightPx; y++ {
			d.Pixels[x][y] = 0
		}
	}
}

// Returns true if any set pixels were changed to unset
// false otherwise
func (d *Display) DrawSprite(x byte, y byte, row byte) bool {
	unset := false

	for i := x; i < 8; i++ {
		xIndex := i % ScreenWidthPx
		wasSet := d.Pixels[xIndex][y] == 0x1
		set := row >> (8 - i) & 0x1
		d.Pixels[xIndex][y] = set
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
