package emu

type Display struct {
	Pixels [ScreenWidthPx * ScreenHeightPx]byte
}

const (
	ScreenWidthPx  = 64
	ScreenHeightPx = 32

	SpriteWidthPx  = 8
	SpriteHeightPx = 8
)

func (d *Display) Draw(mem *Memory) {
	if !mem.ShouldDraw() {
		return
	}
}
