package emu

type CPU struct {
	V     [16]byte   //data registers V0-VF
	I     uint16     //index (memory address) register
	PC    uint16     //program counter
	SP    byte       //stack pointer
	Stack [16]uint16 //stack

	DelayTimer byte //used for timing game events, can be set/read
	SoundTimer byte //beeps when value is nonzero
}

func (c *CPU) Setup() {
	c.PC = RamProgramStart
}
