package emu

import "testing"

// Fx65 - LD Vx, [I]
func TestOpFx65(t *testing.T) {
	e := new(Emulator)
	e.Setup("../roms/IBM.ch8")

	expected := byte(0xAA)
	address := 0x400
	e.cpu.I = uint16(address)
	e.memory.RAM[address] = expected
	i := uint16(0xFF65)
	e.opFx65(i)

	for i := 0; i <= 0xF; i++ {
		actual := e.cpu.V[0]
		if expected != actual {
			t.Errorf("Expected %X at V%d but was %X", expected, i, actual)
		}
	}
}
