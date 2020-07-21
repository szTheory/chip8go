package emu

type CPU struct {
	V  [16]byte //data registers V0-VF
	I  uint16   //index (memory address) register
	PC uint16   //program counter
	SP byte     //stack pointer

	//sound registers
	//when nonzero they decrement at 60 hertz until they reach 0
	DelayTimer byte //used for timing game events, can be set/read
	SoundTimer byte //beep when value is nonzero
}

func (c *CPU) Setup() {
	c.PC = RamProgramStart
}
