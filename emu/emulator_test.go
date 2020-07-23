package emu

import (
	"testing"
)

// Fx33 - LD B, Vx
func TestOpFx33(t *testing.T) {
	e := new(Emulator)
	e.Setup("../roms/IBM.ch8")

	tests := [16][5]byte{
		{0, 255, 2, 5, 5},
	}

	e.cpu.I = RamProgramStart

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		register := test[0]
		decimal := test[1]
		hundreds := test[2]
		tens := test[3]
		ones := test[4]

		e.cpu.V[register] = decimal
		e.opFx33(register)

		if e.memory.RAM[e.cpu.I] != hundreds {
			t.Errorf("Expected hundreds of %d but was %d", hundreds, e.memory.RAM[e.cpu.I])
		}
		if e.memory.RAM[e.cpu.I+1] != tens {
			t.Errorf("Expected tens of %d but was %d", tens, e.memory.RAM[e.cpu.I+1])
		}
		if e.memory.RAM[e.cpu.I+2] != ones {
			t.Errorf("Expected ones of %d but was %d", ones, e.memory.RAM[e.cpu.I+2])
		}
	}
}

// Fx65 - LD Vx, [I]
func TestOpFx65(t *testing.T) {
	e := new(Emulator)
	e.Setup("../roms/IBM.ch8")

	expected := byte(0xAA)
	e.cpu.I = RamProgramStart
	e.memory.RAM[e.cpu.I] = expected
	i := uint16(0xFF65)
	e.opFx65(byte(i & 0xF00 >> 8))

	for i := 0; i <= 0xF; i++ {
		actual := e.cpu.V[0]
		if expected != actual {
			t.Errorf("Expected %X at V%d but was %X", expected, i, actual)
		}
	}
}
