package emu

type Input struct {
	//input is done with a hex keyboard that has 16 keys ranging 0 to F
	keys            [16]bool
	WaitingForInput bool
}

func (i *Input) IsPressed(keyIndex byte) bool {
	return i.keys[keyIndex]
}

func (i *Input) Update(keyIndex byte, pressed bool) {
	i.keys[keyIndex] = pressed
}
