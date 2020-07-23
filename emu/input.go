package emu

type Input struct {
	//input is done with a hex keyboard that has 16 keys ranging 0 to F
	keys [16]bool
}

func (i *Input) IsPressed(keyIndex byte) bool {
	return i.keys[keyIndex]
}

func (i *Input) Press(keyIndex byte) {
	i.keys[keyIndex] = true
}

func (i *Input) Release(keyIndex byte) {
	i.keys[keyIndex] = false
}
